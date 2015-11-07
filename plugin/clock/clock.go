package clock

import (
	"github.com/XANi/uberstatus/uber"
//	"gopkg.in/yaml.v1"
	"time"
	"github.com/op/go-logging"
//	"fmt"
)

var log = logging.MustGetLogger("main")


type Config struct {
	long_format string
	short_format string
	interval int
}


func New(config map[string]interface{}, events chan uber.Event, update chan uber.Update) {
	c := loadConfig(config)
	for {
		select {
		case _ = (<-events):
			Update(update, c.long_format)
		case <-time.After(time.Duration(c.interval) * time.Millisecond):
			Update(update, c.long_format)
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
				c.long_format=converted
			case key == `short_format`:
				c.short_format=converted
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


func Update(update chan uber.Update, format string) {
	time := GetTimeEvent(format)
	time.Color=`#DDDDFF`
	update <- time
}


func GetTimeEvent(format string) uber.Update {
	t :=  time.Now().Local()
	var ev uber.Update
	ev.FullText = t.Format(format)
	ev.ShortText = t.Format(`15:04:05`)
	return ev
}
