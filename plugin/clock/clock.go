package clock

import (
	"github.com/XANi/uberstatus/uber"
	//	"gopkg.in/yaml.v1"
	"github.com/op/go-logging"
	"time"
	//	"fmt"
)

var log = logging.MustGetLogger("main")

type Config struct {
	long_format  string
	short_format string
	interval     int
}

func Run(cfg uber.PluginConfig) {
	c := loadConfig(cfg.Config)
	for {
		select {
		case ev := <-cfg.Events:
			if ev.Button == 3 {
				UpdateWithMonth(cfg.Update)
			} else {
				UpdateWithDate(cfg.Update)
			}
			select {
			// after next click "normal" (time) handler will fire
			case _ = <-cfg.Events:
			case <-time.After(2 * time.Second):
			}
		case _ = <-cfg.Trigger:
			Update(cfg.Update, c.long_format)
		case <-time.After(time.Duration(c.interval) * time.Millisecond):
			Update(cfg.Update, c.long_format)
		}
	}

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
				log.Warning("unknown config key: [%s]", key)

			}
			log.Warning("t: %s", key)
		} else {
			converted, ok := value.(int)
			if ok {
				switch {
				case key == `interval`:
					c.interval = converted
				default:
					log.Warning("unknown config key: [%s]", key)
				}
			} else {
				log.Error("Cant interpret value of config key [%s]", key)
			}
		}
	}
	return c
}

func UpdateWithDate(update chan uber.Update) {
	Update(update, "2006-01-02")
}

func UpdateWithMonth(update chan uber.Update) {
	Update(update, "Mon Jan MST")
}

func Update(update chan uber.Update, format string) {
	time := GetTimeEvent(format)
	time.Color = `#DDDDFF`
	update <- time
}

func GetTimeEvent(format string) uber.Update {
	t := time.Now().Local()
	var ev uber.Update
	ev.FullText = t.Format(format)
	ev.ShortText = t.Format(`15:04:05`)
	return ev
}
