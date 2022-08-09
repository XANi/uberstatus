package i3blocks

import (
	//	"gopkg.in/yaml.v1"
	"bytes"
	"fmt"
	"github.com/XANi/uberstatus/config"
	"github.com/op/go-logging"
	"os"
	"os/exec"
	"strconv"
	"strings"
	//
	"github.com/XANi/uberstatus/uber"
)

var log = logging.MustGetLogger("main")

type pluginConfig struct {
	prefix   string
	Command  string
	Interval int
	Color    string
	name     string
	instance string
}

func New(cfg uber.PluginConfig) (uber.Plugin, error) {
	p, err := loadConfig(cfg.Config)
	return &p, err
}
func (c *pluginConfig) Init() error {
	return nil
}
func (c *pluginConfig) GetUpdateInterval() int {
	return c.Interval
}
func (c *pluginConfig) UpdatePeriodic() uber.Update {
	return c.Update(uber.Event{})
}
func (c *pluginConfig) UpdateFromEvent(ev uber.Event) uber.Update {
	return c.Update(ev)
}

func loadConfig(c config.PluginConfig) (pluginConfig, error) {
	var cfg pluginConfig
	cfg.Interval = 1000
	cfg.Color = `#ffffff`
	return cfg, c.GetConfig(&cfg)
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

func (cfg *pluginConfig) Update(ev uber.Event) uber.Update {
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
	cmd := exec.Command(cfg.Command)
	log.Debug(cfg.Command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		log.Errorf("Error running %s: %s", cfg.Command, err)
	}
	s := out.String()

	st := strings.Split(s, "\n")
	st_len := len(st)

	if st_len >= 4 {
		upd.Color = st[2]
	} else {
		upd.Color = cfg.Color
	}
	if st_len == 3 {
		upd.ShortText = st[1]
	}
	if st_len == 2 {
		upd.ShortText = st[0]
	}
	// len of 1 means there was nothing to split, no \n probably means invalid input
	if st_len <= 1 {
		log.Warningf("Command %s returned nothing", cfg.Command)
		return upd
	} else {
		upd.FullText = fmt.Sprint(cfg.prefix, st[0])
	}
	return upd
}
