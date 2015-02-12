package plugin

import (
	plugin_clock "plugin/clock"
	plugin_network "plugin/network"
	"plugin_interface"
	"fmt"
	"github.com/op/go-logging"
)
var log = logging.MustGetLogger("main")
func NewPlugin(
	name string, // Plugin name
	instance string, // Plugin instance
	config *map[string]interface{}, // Plugin config
	update_filtered chan plugin_interface.Update, // Update channel
) (	chan plugin_interface.Event)  {
	events := make(chan plugin_interface.Event, 16)
	update := make(chan plugin_interface.Update,1)
	log.Info("Adding plugin %s, instance %s",name, instance)
	switch {
	case name == `clock`:
		go plugin_clock.New(config, events, update)
	case name == `network`:
		go plugin_network.New(config, events, update)
	case true:
		panic(fmt.Sprintf("no plugin named %s", name))
	}


	go filterUpdate(name, instance, update ,update_filtered)
	return events
}



func filterUpdate(
	name string,
	instance string,
	update chan plugin_interface.Update,
	update_filtered chan plugin_interface.Update ) {
	for {
		ev := <- update
		ev.Name = name
		ev.Instance = instance
		update_filtered <- ev
	}
}
