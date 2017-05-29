package plugin

import (
	"fmt"
	"github.com/op/go-logging"
	"gopkg.in/yaml.v1"
	//
	"github.com/XANi/uberstatus/plugin/clock"
	"github.com/XANi/uberstatus/plugin/cpu"
	"github.com/XANi/uberstatus/plugin/df"
	"github.com/XANi/uberstatus/plugin/example"
	"github.com/XANi/uberstatus/plugin/i3blocks"
	"github.com/XANi/uberstatus/plugin/memory"
	"github.com/XANi/uberstatus/plugin/network"
	"github.com/XANi/uberstatus/plugin/ping"
	// "github.com/XANi/uberstatus/plugin/weather"
	"github.com/XANi/uberstatus/uber"
	"time"
)

var log = logging.MustGetLogger("main")

var plugins = map[string]func(uber.PluginConfig)(uber.Plugin,error){
	"clock":    clock.New,
	"cpu":      cpu.New,
	"df":       df.New,
	"example":  example.New,
	"i3blocks": i3blocks.New,
	"memory":   memory.New,
	"network":  network.New,
	"ping":     ping.New,
	"weather":  example.New,
}


func NewPlugin(
	name string, // Plugin name
	instance string, // Plugin instance
	backend string, // Plugin backend
	config map[string]interface{}, // Plugin config
	update_filtered chan uber.Update, // Update channel
) (uber.Plugin,error) {
	events := make(chan uber.Event, 1)
	update := make(chan uber.Update, 1)
	trigger := make(chan uber.Trigger, 1)
	log.Infof("Adding plugin %s, instance %s", name, instance)
	str, _ := yaml.Marshal(config)
	log.Warning(string(str))
	pluginCfg := uber.PluginConfig{
		Name:     name,
		Instance: instance,
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
	if p, ok := plugins[backend]; ok {
		plugin, err := p(pluginCfg)
		if err != nil {
			return nil, err
		}
		err = plugin.Init()
		if err != nil {
			return nil, err
		}
		go filterUpdate(name, instance, update, update_filtered)
		go run(plugin.GetUpdateInterval(), events, update, trigger, plugin)
		return plugin, nil
	} else {
		return nil, fmt.Errorf("no plugin named %s", backend)
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
