package example

import (
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	//	"gopkg.in/yaml.v1"
	"fmt"
	"github.com/op/go-logging"
	"time"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

var log = logging.MustGetLogger("main")

// set up a config struct
type config struct {
	prefix   string
	interval int
}

type state struct {
	cfg config
	cnt int
	ev  int
}

func Run(cfg uber.PluginConfig) {
	var st state
	st.cfg = loadConfig(cfg.Config)
	_ = fmt.Sprintf("",cfg)
	// initial update on start
	cfg.Update <- st.updatePeriodic()
	for {
		select {
		// call update when user clicked on the plugin
		case updateEvent := (<-cfg.Events):
			cfg.Update <- st.updateFromEvent(updateEvent)
			// that will wait 10 seconds on no event a
			// and it will "eat" next event to switch to "normal" display
			// basically making it "toggle" between two different views
			select {
			case _ = <-cfg.Events:
				cfg.Update <- st.updatePeriodic()
			case <-time.After(10 * time.Second):
			}
		// update on trigger from main code, this can be used to make all widgets update at the same time if that way is preferred over async
		case _ = <-cfg.Trigger:
			cfg.Update <- st.updatePeriodic()
		// update every interval if nothing triggered update before tat
		case <-time.After(time.Duration(st.cfg.interval) * time.Millisecond):
			cfg.Update <- st.updatePeriodic()
		}
	}
}

func (state *state) updatePeriodic() uber.Update {
	var update uber.Update
	// TODO precompile and preallcate
	tpl, _ := util.NewTemplate("uberEvent",`{{color "#00aa00" "Example plugin"}}`)
	update.FullText =  tpl.ExecuteString(nil)
	update.Markup = `pango`
	update.ShortText = `nope`
	update.Color = `#66cc66`
	state.cnt++
	return update
}

func (state *state) updateFromEvent(e uber.Event) uber.Update {
	var update uber.Update
	tpl, _ := util.NewTemplate("uberEvent",`{{printf "%+v" .}}`)
	update.FullText =  tpl.ExecuteString(e)
	update.ShortText = `upd`
	update.Color = `#cccc66`
	state.ev++
	return update
}

// parse received structure into config
func loadConfig(c map[string]interface{}) config {
	var cfg config
	cfg.interval = 10000
	cfg.prefix = "ex: "
	for key, value := range c {
		converted, ok := value.(string)
		if ok {
			switch {
			case key == `prefix`:
				cfg.prefix = converted
			default:
				log.Warningf("unknown config key: [%s]", key)

			}
		} else {
			converted, ok := value.(int)
			if ok {
				switch {
				case key == `interval`:
					cfg.interval = converted
				default:
					log.Warningf("unknown config key: [%s]", key)
				}
			} else {
				log.Errorf("Cant interpret value of config key [%s]", key)
			}
		}
	}
	return cfg
}
