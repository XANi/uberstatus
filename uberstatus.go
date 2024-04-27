package main

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
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

var log *zap.SugaredLogger
var debug = flag.Bool("d", false, "enable debug server on port 6060[pprof]")
var configFile = flag.String("config", "", "path to config file")

var version string

func init() {
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	// naive systemd detection. Drop timestamp if running under it
	if os.Getenv("INVOCATION_ID") != "" || os.Getenv("JOURNAL_STREAM") != "" {
		consoleEncoderConfig.TimeKey = ""
	}
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
	consoleStderr := zapcore.Lock(os.Stderr)
	_ = consoleStderr
	logLevel := zapcore.InfoLevel
	if *debug {
		logLevel = zapcore.DebugLevel
	} else {
	}

	// if needed point differnt priority log to different place
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if *debug {
			return lvl < logLevel
		} else {
			return false
		}
	})
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, os.Stderr, lowPriority),
		zapcore.NewCore(consoleEncoder, os.Stderr, highPriority),
	)
	logger := zap.New(core)
	if *debug {
		logger = logger.WithOptions(
			zap.Development(),
			zap.AddCaller(),
			zap.AddStacktrace(highPriority),
		)
	} else {
		logger = logger.WithOptions(
			zap.AddCaller(),
		)
	}
	log = logger.Sugar()

}

type Config struct {
	Plugins *map[string]map[string]interface{}
}

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
		pluginCfg.Logger = log.Named(pluginCfg.Name)
		log.Infof("Loading plugin %s into slot %d: %+v", pluginCfg.Plugin, idx, pluginCfg)
		if plugins.slotMap[pluginCfg.Name] == nil {
			plugins.slotMap[pluginCfg.Name] = make(map[string]int)
			plugins.plugins[pluginCfg.Name] = make(map[string]uber.Plugin)
		}

		plugins.slotMap[pluginCfg.Name][pluginCfg.Instance] = idx
		plugins.slots[idx] = i3bar.NewMsg()
		plugins.plugins[pluginCfg.Name][pluginCfg.Instance], err = plugin.NewPlugin(pluginCfg, updates)
		if err != nil {
			plugins.slots[idx].FullText = fmt.Sprintf("can't init plugin instance %s: %s", pluginCfg.Instance, err)
			plugins.slots[idx].ShortText = fmt.Sprintf("%s[%s] failed", pluginCfg.Instance, pluginCfg.Name)
			if cfg.PanicOnBadPlugin {
				log.Panicf("Can't initialize plugin [%s:%s]: [%+v]", pluginCfg.Name, pluginCfg.Instance, err)
			}
			pluginCfg.Logger.Errorf("Can't initialize plugin [%s:%s]: [%+v]", pluginCfg.Name, pluginCfg.Instance, err)

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
			//log.Info("Time passed")
		}
		out := `[`
		for idx, msg := range plugins.slots {
			out = out + string(msg.Encode())
			if idx+1 < (len(plugins.slots)) {
				out = out + `,`
			}
		}
		out = out + "],\n"
		select {
		case ow <- out:
		default:
			log.Warn("output channel full, discarding output!")
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
		log.Warnf("Got msg from unknown place, name: %s, instance: %s", update.Name, update.Instance)
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
