package i3blocks

import (
	//	"gopkg.in/yaml.v1"
	"bytes"
	"fmt"
	"github.com/op/go-logging"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	//
	"github.com/XANi/uberstatus/uber"
)

var log = logging.MustGetLogger("main")

type config struct {
	prefix   string
	command  string
	interval int
	color    string
	name     string
	instance string
}

func Run(cfg uber.PluginConfig) {
	c := loadConfig(cfg.Config)
	var nullEv uber.Event
	for {
		select {
		case ev := (<-cfg.Events):
			c.Update(cfg.Update, c, ev)
		case _ = <-cfg.Trigger:
			c.Update(cfg.Update, c, nullEv)
		case <-time.After(time.Duration(c.interval) * time.Millisecond):
			c.Update(cfg.Update, c, nullEv)
		}
	}

}

func loadConfig(raw map[string]interface{}) config {
	var c config
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

//   BLOCK_NAME
//              The name of the block (usually the section name).
//
//       BLOCK_INSTANCE
//              An optional argument to the script.
//
//       BLOCK_BUTTON
//              Mouse button (1, 2 or 3) if the block was clicked.
//
//       BLOCK_X and BLOCK_Y
//              Coordinates where the click occurred, if the block was clicked.

func (c *config) Update(update chan uber.Update, cfg config, ev uber.Event) {
	var upd uber.Update

	os.Setenv("BLOCK_NAME", c.name)
	os.Setenv("BLOCK_NAME", c.instance)
	// no event
	if ev.Button == 0 {
		os.Unsetenv("BLOCK_BUTTON")
		os.Unsetenv("BLOCK_X")
		os.Unsetenv("BLOCK_Y")
	} else {
		os.Setenv("BLOCK_BUTTON", strconv.Itoa(ev.Button))
		os.Setenv("BLOCK_X", strconv.Itoa(ev.X))
		os.Setenv("BLOCK_Y", strconv.Itoa(ev.Y))
	}
	cmd := exec.Command(cfg.command)
	log.Debug(cfg.command)
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
		upd.Color = st[2]
	} else {
		upd.Color = cfg.color
	}
	if st_len == 3 {
		upd.ShortText = st[1]
	}
	if st_len == 2 {
		upd.ShortText = st[0]
	}
	// len of 1 means there was nothing to split, no \n probably means invalid input
	if st_len <= 1 {
		log.Warning("Command %s returned nothing", cfg.command)
		return
	} else {
		upd.FullText = fmt.Sprint(cfg.prefix, st[0])
	}
	update <- upd
}
