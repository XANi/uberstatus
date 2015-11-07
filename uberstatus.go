package main

import (
	"encoding/json"
	"fmt"
	"github.com/op/go-logging"
	"gopkg.in/yaml.v1"
	//	"io/ioutil"
	"os"
	"regexp"
	"time"
	//
	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/i3bar"
	"github.com/XANi/uberstatus/plugin"
	"github.com/XANi/uberstatus/uber"
)

type Config struct {
	Plugins *map[string]map[string]interface{}
}

var log = logging.MustGetLogger("main")
var logFormat = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortpkg}↛%{shortfunc}: %{level:.4s} %{id:03x} %{color:reset}%{message}",
)

func main() {
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logBackendFormatter := logging.NewBackendFormatter(logBackend, logFormat)
	_ = logBackendFormatter
	logBackendLeveled := logging.AddModuleLevel(logBackendFormatter)
	logBackendLeveled.SetLevel(logging.NOTICE, "")
	logging.SetBackend(logBackendLeveled)

	log.Info("Starting")
	header := i3bar.NewHeader()
	msg := i3bar.NewMsg()
	msg.FullText = `test`
	b, err := json.Marshal(header)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	//c, err := json.Marshal(msg)

	i3input := i3bar.EventReader()
	updates := make(chan uber.Update, 10)
	cfg := config.LoadConfig()
	slotMap := make(map[string]map[string]int)
	slots := make([]i3bar.Msg, len(cfg.Plugins))
	for idx, pluginCfg := range cfg.Plugins {
		log.Info("Loading plugin %s into slot %d", pluginCfg.Plugin, idx)
		if slotMap[pluginCfg.Plugin] == nil {
			slotMap[pluginCfg.Plugin] = make(map[string]int)
		}
		if len(pluginCfg.Instance) == 0 {
			pluginCfg.Instance = pluginCfg.Name
		}

		slotMap[pluginCfg.Plugin][pluginCfg.Instance] = idx
		slots[idx] = i3bar.NewMsg()
		_ = plugin.NewPlugin(pluginCfg.Plugin, pluginCfg.Instance, pluginCfg.Config, updates)
	}

	_ = slots
	// fmt.Println("\n[")

	// plugins := config.Plugins
	// ifd := (*plugins)[`clock`] //.(map[string]interface{})
	// net := (*plugins)[`clock`] //.(map[string]interface{})
	// //	_ = plugin.NewPlugin("clock", "", &ifd, updates)
	// _ = plugin.NewPlugin("network", "", &net, updates)
	// _ = ifd
	fmt.Println(`[`)

	for {
		select {
		case ev := (<-i3input):
			log.Info("Gut event from plugin %d", ev.Button)
		case upd := <-updates:
			parseUpdate(upd, &slots, &slotMap)
		case <-time.After(time.Second):
			log.Info("Time passed")
		}
		fmt.Print(`[`)

		for idx, msg := range slots {
			os.Stdout.Write(msg.Encode())
			if idx+1 < (len(slots)) {
				fmt.Print(`,`)
			}
		}
		fmt.Println(`],`)

	}
}

func parseUpdate(update uber.Update, slots *[]i3bar.Msg, slotMap *map[string]map[string]int) {
	if val, ok := (*slotMap)[update.Name][update.Instance]; ok {
		(*slots)[val] = i3bar.CreateMsg(update)
		log.Info("Got msg from unk name: %s, instance: %s", update.Name, update.Instance)
	} else {
		log.Warning("Got msg from unknown place, name: %s, instance: %s", update.Name, update.Instance)
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
