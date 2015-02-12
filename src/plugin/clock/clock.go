package clock

import (
	"plugin_interface"
//	"gopkg.in/yaml.v1"
	"time"
//	"fmt"
)

type Config struct {
	long_format string
	short_format string
}


func New(config map[string]interface{}, events chan plugin_interface.Event, update chan plugin_interface.Update) {
	c := loadConfig(config)
	for {
		select {
		case _ = (<-events):
			Update(update, c.long_format)
		case <-time.After(time.Second):
			Update(update, c.long_format)
		}
	}

}

func loadConfig(raw map[string]interface{}) Config {
	var c Config
	c.long_format = `2006-01-02 MST 15:04:05.00`
	c.short_format = `15:04:05`
	for key, value := range raw {
		converted, ok := value.(string)
		if ok {
			switch {
			case key == `long_format`:
				c.long_format=converted

			case key == `short_format`:
				c.short_format=converted
			}
			} else {
			_ = ok
		}
	}
	return c
}


func Update(update chan plugin_interface.Update, format string) {
	time := GetTimeEvent(format)
	time.Color=`#DDDDFF`
	update <- time
}


func GetTimeEvent(format string) plugin_interface.Update {
	t :=  time.Now().Local()
	var ev plugin_interface.Update
	ev.FullText = t.Format(format)
	ev.ShortText = t.Format(`15:04:05`)
	return ev
}
