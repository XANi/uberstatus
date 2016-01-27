package plugin

import (
	"fmt"
	"github.com/op/go-logging"
	"gopkg.in/yaml.v1"
	//
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/plugin/clock"
	"github.com/XANi/uberstatus/plugin/cpu"
	"github.com/XANi/uberstatus/plugin/network"
	"github.com/XANi/uberstatus/plugin/i3blocks"
	"github.com/XANi/uberstatus/plugin/example"

)
var log = logging.MustGetLogger("main")


var plugins =  map[string]func(uber.PluginConfig)  {
	"clock": clock.Run,
	"cpu": cpu.Run,
	"network": network.Run,
	"i3blocks": i3blocks.Run,
	"example": example.Run,
};


func NewPlugin(
	name string, // Plugin name
	instance string, // Plugin instance
	backend string, // Plugin backend
	config map[string]interface{}, // Plugin config
	update_filtered chan uber.Update, // Update channel
) (	chan uber.Event)  {
	events := make(chan uber.Event, 1)
	update := make(chan uber.Update,1)
	log.Info("Adding plugin %s, instance %s",name, instance)
	str, _ := yaml.Marshal(config)
	log.Warning(string(str))
	if p, ok := plugins[backend]; ok {
		go p(uber.PluginConfig{
			Name: name,
			Instance: instance,
			Config: config,
			Events: events,
			Update: update,
		})
	} else {
		log.Error(fmt.Sprintf("no plugin named %s", backend))
		panic(fmt.Sprintf("no plugin named %s", backend))
	}
	go filterUpdate(name, instance, update ,update_filtered)
	return events
}



func filterUpdate(
	name string,
	instance string,
	update chan uber.Update,
	update_filtered chan uber.Update ) {
	for {
		ev := <- update
		ev.Name = name
		ev.Instance = instance
		update_filtered <- ev
	}
}
