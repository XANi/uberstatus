package i3blocks

import (
	//	"gopkg.in/yaml.v1"
	"os/exec"
	"time"
	"github.com/op/go-logging"
	"strings"
	"bytes"
	"fmt"
	//
	"github.com/XANi/uberstatus/uber"

)

var log = logging.MustGetLogger("main")


type Config struct {
	prefix string
	command string
	interval int
	color string
}

func New(config map[string]interface{}, events chan uber.Event, update chan uber.Update) {
	c := loadConfig(config)
	for {
		select {
		case _ = (<-events):
			Update(update,c)
		case <-time.After(time.Duration(c.interval) * time.Millisecond):
			Update(update,c)
		}
	}

}

func loadConfig(raw map[string]interface{}) Config {
	var c Config
	c.interval = 1000
	c.color = `#ffffff`
	log.Info("loading config")
	for key, value := range raw {
		converted, ok := value.(string)
		if ok {
			switch {
			case key == `prefix`:
				c.prefix = converted
			case key == `command`:
				c.command = converted
			case key == `color`:
				c.color = converted
			default:
				log.Warning("unknown config key: [%s]", key)

			}
			log.Warning("t: %s %s", key, converted)
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


func Update(update chan uber.Update, cfg Config) {
	var ev uber.Update
	cmd := exec.Command(cfg.command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
	s := out.String()

    st := strings.Split(s, "\n")
	st_len := len(st)

	if st_len >= 4 {
		ev.Color = st[2]
	} else {
		ev.Color = cfg.color
	}
	if st_len == 3 {
		ev.ShortText = st[1]
	}
	if st_len == 2 {
		ev.ShortText = st[0]
	}
	// len of 1 means there was nothing to split, no \n probably means invalid input
	if st_len <= 1 {
		log.Warning("Command %s returned nothing",cfg.command)
		return
	} else {
		ev.FullText = fmt.Sprint(cfg.prefix, st[0])
	}
	update <- ev
}
