package main

import (
	"encoding/json"
	"fmt"
	"github.com/op/go-logging"
	"gopkg.in/yaml.v1"
	//	"runtime"
	//	"io/ioutil"
	"flag"
	"net/http"
	_ "net/http/pprof"
	"os"
	"regexp"
	"runtime"
	"time"
	//
	"github.com/XANi/go-yamlcfg"
	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/i3bar"
	"github.com/XANi/uberstatus/plugin"
	"github.com/XANi/uberstatus/uber"
)

type Config struct {
	Plugins *map[string]map[string]interface{}
}

var debug = flag.Bool("d", false, "enable debug server on port 6060[pprof]")
var configFile = flag.String("config", "", "path to config file")

var version string
var log = logging.MustGetLogger("main")
var logFormat = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortpkg}↛%{shortfunc}: %{level:.4s} %{id:03x} %{color:reset}%{message}",
)

type pluginMap struct {
	// channels used to send events to plugin
	plugins map[string]map[string]uber.Plugin
	slots   []i3bar.Msg
	slotMap map[string]map[string]int
}

func main() {
	flag.Parse()
	if *debug {
		go func() {
			runtime.SetCPUProfileRate(10000)
			log.Errorf("%+v", http.ListenAndServe("127.0.0.1:6060", nil))
		}()
	}
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logBackendFormatter := logging.NewBackendFormatter(logBackend, logFormat)
	_ = logBackendFormatter
	logBackendLeveled := logging.AddModuleLevel(logBackendFormatter)
	logBackendLeveled.SetLevel(logging.NOTICE, "")
	logging.SetBackend(logBackendLeveled)
	var cfg config.Config
	var err error
	if len(*configFile) > 0 {
		err = yamlcfg.LoadConfig([]string{*configFile}, &cfg)
	} else {
		err = yamlcfg.LoadConfig(config.CfgFiles, &cfg)
	}
	if err != nil {
		log.Errorf("Can't load config: %s", err)
		os.Exit(1)
	}
	log.Debug("config: %+v", &cfg)
	log.Info("Starting")
	header := i3bar.NewHeader()
	msg := i3bar.NewMsg()
	msg.FullText = `test`
	b, err := json.Marshal(header)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)

	i3input := i3bar.EventReader()
	// FIXME
	// it looks like i3bar blocks stdout when stdin is not handled
	// which feedbacks into plugins if channel does
	updates := make(chan uber.Update, 100)
	plugins := pluginMap{
		slotMap: make(map[string]map[string]int),
		slots:   make([]i3bar.Msg, len(cfg.Plugins)),
		plugins: make(map[string]map[string]uber.Plugin),
	}
	for idx, pluginCfg := range cfg.Plugins {
		log.Infof("Loading plugin %s into slot %d: %+v", pluginCfg.Plugin, idx, pluginCfg)
		if plugins.slotMap[pluginCfg.Name] == nil {
			plugins.slotMap[pluginCfg.Name] = make(map[string]int)
			plugins.plugins[pluginCfg.Name] = make(map[string]uber.Plugin)
		}

		plugins.slotMap[pluginCfg.Name][pluginCfg.Instance] = idx
		plugins.slots[idx] = i3bar.NewMsg()
		plugins.plugins[pluginCfg.Name][pluginCfg.Instance], err = plugin.NewPlugin(pluginCfg.Name, pluginCfg.Instance, pluginCfg.Plugin, pluginCfg.Config, updates)
		if err != nil {
			log.Panicf("Can't initialize plugin: [%+v]", err)
		}
	}

	fmt.Println(`[`)
	ow := make(chan string, 32)
	go outputWriter(ow)
	for {
		select {
		case ev := (<-i3input):
			plugins.parseEvent(ev)
		case upd := <-updates:
			plugins.parseUpdate(upd)
		case <-time.After(time.Second * 1):
			log.Info("Time passed")
		}
		out := `[`
		for idx, msg := range plugins.slots {
			out = out + string(msg.Encode())
			if idx+1 < (len(plugins.slots)) {
				out = out + `,`
			}
		}
		out = out + `],`
		select {
		case ow <- out:
		default:
			log.Warning("output channel full, discarding output!")
		}
	}
}
func outputWriter(out chan string) {
	for {
		ev := <-out
		os.Stdout.Write([]byte(ev))
	}
}

func (plugins *pluginMap) parseUpdate(update uber.Update) {
	if val, ok := plugins.slotMap[update.Name][update.Instance]; ok {
		plugins.slots[val] = i3bar.CreateMsg(update)
	} else {
		log.Warningf("Got msg from unknown place, name: %s, instance: %s", update.Name, update.Instance)
	}
}

func (plugins *pluginMap) parseEvent(ev uber.Event) {
	if val, ok := plugins.plugins[ev.Name][ev.Instance]; ok {
		upd := val.UpdateFromEvent(ev)
		upd.Name = ev.Name
		upd.Instance = ev.Instance
		plugins.parseUpdate(upd)
	} else {
		log.Infof("rejected event %+v", ev)
	}

}

func getTime() []byte {
	msg := i3bar.NewMsg()
	msg.Name = "clock"
	t := time.Now().Local()
	// reference Mon Jan 2 15:04:05 MST 2006 (unix: 1136239445)
	msg.FullText = t.Format(`15:04:05`)
	msg.Color = `#ffffff`
	return msg.Encode()
}

func San(in []byte) []byte {
	re := regexp.MustCompile(`\,{`)
	return re.ReplaceAllLiteral(in, []byte(`{`))
}

func PrintInterface(a interface{}) {
	fmt.Println("Interface:")
	txt, _ := yaml.Marshal(a)
	fmt.Printf("%s", txt)
}
