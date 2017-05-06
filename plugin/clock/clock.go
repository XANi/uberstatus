package clock

import (
	"github.com/XANi/uberstatus/uber"
	//	"gopkg.in/yaml.v1"
	"github.com/op/go-logging"
	"time"
	//	"fmt"
)

var log = logging.MustGetLogger("main")

type plugin struct {
	cfg Config
	nextTs time.Time
}

type Config struct {
	long_format  string
	short_format string
	interval     int
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
	t := time.Now()
	var	ev uber.Update
	// sleep if we are still displaying event data
	for p.nextTs.After(t) {
		diff :=p.nextTs.Sub(t)
		// cap sleeping at 10s in case date changes between ticks
		if diff > time.Second * 10  {
			//time.Sleep(time.Second * 10)
			time.Sleep(diff)
		} else {
			time.Sleep(diff)
		}
		t = time.Now()
	}
	ev = p.GetTimeEvent(t.Local(), p.cfg.long_format)
	t = time.Now().Local()
	//ev := GetTimeEvent(t, p.cfg.long_format)
	ev.Color = `#DDDDFF`
	return ev
}
func (p *plugin) UpdateFromEvent(ev uber.Event) uber.Update {
	var upd uber.Update
	t := time.Now()
	if ev.Button == 3 {
		upd = p.GetTimeEvent(t.Local(),"2006-01-02")
	} else {
		upd = p.GetTimeEvent(t.Local(),"Mon Jan MST")
	}
	p.nextTs = t.Add(time.Second * 3)
	return upd
}

func loadConfig(raw map[string]interface{}) Config {
	var c Config
	c.long_format = `2006-01-02 MST 15:04:05.00`
	c.short_format = `15:04:05`
	c.interval = 500
	for key, value := range raw {
		converted, ok := value.(string)
		if ok {
			switch {
			case key == `long_format`:
				c.long_format = converted
			case key == `short_format`:
				c.short_format = converted
			default:
				log.Warningf("unknown config key: [%s]", key)

			}
			log.Warningf("t: %s", key)
		} else {
			converted, ok := value.(int)
			if ok {
				switch {
				case key == `interval`:
					c.interval = converted
				default:
					log.Warningf("unknown config key: [%s]", key)
				}
			} else {
				log.Errorf("Cant interpret value of config key [%s]", key)
			}
		}
	}
	return c
}

func (p *plugin) GetTimeEvent(t time.Time, format string) uber.Update {
	var ev uber.Update
	ev.FullText = t.Format(format)
	ev.ShortText = t.Format(p.cfg.short_format)
	return ev
}
