package plugin

import (
	"fmt"
	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/plugin/gpu"
	"github.com/XANi/uberstatus/plugin/mqtt"
	"github.com/XANi/uberstatus/plugin/syncthing"

	"github.com/op/go-logging"
	"gopkg.in/yaml.v1"
	//
	"github.com/XANi/uberstatus/plugin/clock"
	"github.com/XANi/uberstatus/plugin/cpu"
	"github.com/XANi/uberstatus/plugin/cpufreq"
	"github.com/XANi/uberstatus/plugin/example"
	"github.com/XANi/uberstatus/plugin/i3blocks"
	"github.com/XANi/uberstatus/plugin/memory"
	"github.com/XANi/uberstatus/plugin/network"
	"github.com/XANi/uberstatus/plugin/ping"
	"github.com/XANi/uberstatus/plugin/pomodoro"
        "github.com/XANi/uberstatus/plugin/uptime"
	"github.com/XANi/uberstatus/plugin/weather"
	"github.com/XANi/uberstatus/uber"
	"time"
	"github.com/XANi/uberstatus/plugin/df"
	"github.com/XANi/uberstatus/plugin/debug"
	"github.com/XANi/uberstatus/plugin/pipe"
)

var log = logging.MustGetLogger("main")

var plugins = map[string]func(uber.PluginConfig)(uber.Plugin,error){
	"clock":     clock.New,
	"cpu":       cpu.New,
	"cpufreq":   cpufreq.New,
	"debug":     debug.New,
	"df":        df.New,
	"example":   example.New,
	"gpu":       gpu.New,
	"i3blocks":  i3blocks.New,
	"memory":    memory.New,
	"mqtt":      mqtt.New,
	"network":   network.New,
	"ping":      ping.New,
	"pipe":      pipe.New,
	"pomodoro":  pomodoro.New,
	"syncthing": syncthing.New,
	"uptime":    uptime.New,
	"weather":   weather.New,
}


func NewPlugin(
	config config.PluginConfig,
	update_filtered chan uber.Update, // Update channel
) (uber.Plugin,error) {
	events := make(chan uber.Event, 1)
	update := make(chan uber.Update, 1)
	trigger := make(chan uber.Trigger, 1)
	log.Infof("Adding plugin %s, instance %s", config.Name, config.Instance)
	str, _ := yaml.Marshal(config)
	log.Debug(string(str))
	pluginCfg := uber.PluginConfig{
		Config:   config,
		Update:   update,
	}
	// TODO make it global somehow
	// interval := 1000
	// if val, ok := config["interval"]; ok {
	// 	if ok {
	// 		converted, ok := val.(int)
	// 		if ok {
	// 			interval = converted
	// 		}
	// 	}
	// }
	if p, ok := plugins[config.Plugin]; ok {
		plugin, err := p(pluginCfg)
		if err != nil {
			return nil, err
		}
		err = plugin.Init()
		if err != nil {
			return nil, err
		}
		go filterUpdate(config.Name, config.Instance, update, update_filtered)
		go run(plugin.GetUpdateInterval(), events, update, trigger, plugin)
		return plugin, nil
	} else {
		return nil, fmt.Errorf("no plugin named %s", config.Plugin)
	}
}

func run(interval int, events chan uber.Event, update chan uber.Update, trigger chan uber.Trigger, p uber.Plugin) {
	//initial update so we have something to display when main loop runs first time
	update <- p.UpdatePeriodic()
	// run periodic updates independenely of on-demand ones
	go func() {
		for {
			if interval > 0 {
				select {
				case _ = <-trigger:
					update <- p.UpdatePeriodic()
				case <-time.After(time.Duration(interval) * time.Millisecond):
					update <- p.UpdatePeriodic()
				}
			} else {
				select {
				case _ = <-trigger:
					update <- p.UpdatePeriodic()
				}
			}
		}
	}()

	for {
		select {
		case updateEv := <-events:
			update <- p.UpdateFromEvent(updateEv)
		}
	}

}

func filterUpdate(
	name string,
	instance string,
	update chan uber.Update,
	update_filtered chan uber.Update) {
	for {
		ev := <-update
		ev.Name = name
		ev.Instance = instance
		update_filtered <- ev
	}
}
