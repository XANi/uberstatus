package cpu

import (
	"plugin_interface"
//	"gopkg.in/yaml.v1"
	"github.com/op/go-logging"
//	"fmt"
)

var log = logging.MustGetLogger("main")


type Config struct {
	interval int
}
type CpuStats struct {
	user uint64
	nice uint64
	system uint64
	idle uint64
	iowait uint64
	irq uint64
	softirq uint64
	steal uint64
	guest uint64
	guest_nice uint64
	ts time.Time
}


func New(config map[string]interface{}, events chan plugin_interface.Event, update chan plugin_interface.Update) {
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
	c.interval = 500
	for key, value := range raw {
		converted, ok := value.(string)
		if ok {
			switch {
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
