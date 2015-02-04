package plugin

import (
	plugin_clock "plugin/clock"
	"plugin_interface"
	"fmt"
)
func NewPlugin(name string, config interface{}) (chan plugin_interface.Event, chan plugin_interface.Update)  {
	events := make(chan plugin_interface.Event, 16)
	update := make(chan plugin_interface.Update)

	switch {
	case name == `clock`:
		go plugin_clock.New(config, events,update)
	case true:
		panic(fmt.Sprintf("no plugin named %s", name))
	}
	return events, update
}
