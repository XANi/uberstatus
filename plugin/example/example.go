package example

import (
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	//	"gopkg.in/yaml.v1"
	"github.com/op/go-logging"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

var log = logging.MustGetLogger("main")

// set up a config struct
type config struct {
	prefix   string
	interval int
}

type plugin struct {
	cfg config
	cnt int
	ev  int
}

func New(cfg uber.PluginConfig) (uber.Plugin, error) {
	p := &plugin{}
	p.cfg = loadConfig(cfg.Config)
	return  p, nil
}

func (p *plugin) Init() error {
	return nil
}

func (p *plugin) UpdatePeriodic() uber.Update {
	var update uber.Update
	// TODO precompile and preallcate
	tpl, _ := util.NewTemplate("uberEvent",`{{color "#00aa00" "Example plugin"}}{{.}}`)
	update.FullText =  tpl.ExecuteString(p.cnt)
	update.Markup = `pango`
	update.ShortText = `nope`
	update.Color = `#66cc66`
	p.cnt++
	return update
}

func (p *plugin) UpdateFromEvent(e uber.Event) uber.Update {
	var update uber.Update
	tpl, _ := util.NewTemplate("uberEvent",`{{printf "%+v" .}}`)
	update.FullText =  tpl.ExecuteString(e)
	update.ShortText = `upd`
	update.Color = `#cccc66`
	p.ev++
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
