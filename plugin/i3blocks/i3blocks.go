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
func New(cfg uber.PluginConfig) (uber.Plugin, error) {
	p := loadConfig(cfg.Config)
	return  &p, nil
}
func (c *config) Init() error {
	return nil
}

func (c *config) UpdatePeriodic() uber.Update {
	return c.Update(uber.Event{})
}
func (c *config) UpdateFromEvent(ev uber.Event) uber.Update {
	return c.Update(ev)
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
				log.Warningf("unknown config key: [%s]", key)

			}
			log.Warningf("t: %s %s", key, converted)
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

func (cfg *config) Update(ev uber.Event) uber.Update {
	var upd uber.Update

	os.Setenv("BLOCK_NAME", cfg.name)
	os.Setenv("BLOCK_NAME", cfg.instance)
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
		log.Warningf("Command %s returned nothing", cfg.command)
		return upd
	} else {
		upd.FullText = fmt.Sprint(cfg.prefix, st[0])
	}
	return upd
}
