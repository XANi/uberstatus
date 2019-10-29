package clock

import (
	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	//	"gopkg.in/yaml.v1"
	"github.com/op/go-logging"
	"time"
	//	"fmt"
)

var log = logging.MustGetLogger("main")

type plugin struct {
	cfg pluginConfig
	nextTs time.Time
}

type pluginConfig struct {
	Long_format  string
	Short_format string
	Interval     int
}

func New(cfg uber.PluginConfig) (z uber.Plugin, err error) {
	p := &plugin{}
	p.cfg, err = loadConfig(cfg.Config)
	return  p, nil
}


func (p *plugin) Init() error {
	return nil
}

func (p *plugin) GetUpdateInterval() int {
	return p.cfg.Interval
}

func (p *plugin) UpdatePeriodic() uber.Update {
	t := time.Now()
	var	ev uber.Update
	util.WaitForTs(&p.nextTs)
	ev = p.GetTimeEvent(t.Local(), p.cfg.Long_format)
	t = time.Now().Local()
	//ev := GetTimeEvent(t, p.cfg.Long_format)
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

func loadConfig(c config.PluginConfig) (pluginConfig ,error) {
	var cfg pluginConfig
	cfg.Long_format = `2006-01-02 MST 15:04:05.00`
	cfg.Short_format = `15:04:05`
	cfg.Interval = 500
	return cfg, c.GetConfig(&cfg)
}

func (p *plugin) GetTimeEvent(t time.Time, format string) uber.Update {
	var ev uber.Update
	ev.FullText = t.Format(format)
	ev.ShortText = t.Format(p.cfg.Short_format)
	return ev
}
