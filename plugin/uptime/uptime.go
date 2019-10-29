package uptime

import (
	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	//	"gopkg.in/yaml.v1"
	"github.com/op/go-logging"
	"time"
	"bufio"
	"fmt"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

var log = logging.MustGetLogger("main")

// set up a pluginConfig struct
type pluginConfig struct {
	Prefix   string
	Interval int
}

type plugin struct {
	cfg               pluginConfig
	uptimeReader      uptimeReader
	uptimeFileScanner *bufio.Scanner
	uptimeTpl         *util.Template
	uptimeTplShort    *util.Template
	dynamicInterval   int
	nextTs            time.Time
}

type uptimeTpl struct {
	Prefix string
	Uptime string
	UptimeShort string
}

func New(cfg uber.PluginConfig) (z uber.Plugin, err error) {
	p := &plugin{}
	p.cfg, err = loadConfig(cfg.Config)
	return  p, nil
}

func (p *plugin) Init() (err error) {
	p.uptimeTpl, err = util.NewTemplate("uptime",`{{printf "%s %s" .Prefix .Uptime}}`)
	if err != nil { return }
	p.uptimeTplShort, err = util.NewTemplate("uptimeShort",`u:{{.UptimeShort}}`)
	return
}

func (p *plugin) GetUpdateInterval() int {
	return p.cfg.Interval
}
func (p *plugin) UpdatePeriodic() uber.Update {
	var update uber.Update
	uptime := p.getUptime()
	uptimeValue :=  uptimeTpl {
		Prefix: p.cfg.Prefix,
		Uptime: fmt.Sprintf("%7s", util.FormatDuration(uptime)),
		UptimeShort: util.FormatDuration(uptime),
	}
	update.FullText =  p.uptimeTpl.ExecuteString(&uptimeValue)
	update.ShortText = p.uptimeTplShort.ExecuteString(&uptimeValue)
	update.Markup = `pango`
	update.Color = `#66cc66`
	return update
}

func (p *plugin) UpdateFromEvent(e uber.Event) uber.Update {
	// set next TS updatePeriodic will wait to.
	p.nextTs = time.Now().Add(time.Second * 3)
	var update uber.Update
	uptime := p.getUptime()
	ts := time.Now().Add(uptime * -1)
	update.FullText = ts.Format(`2006-01-02 MST 15:04:05.00`)
	update.ShortText = ts.Format(`15:04:05`)
	return update
}

// parse received structure into pluginConfig
func loadConfig(c config.PluginConfig) (pluginConfig,error) {
	var cfg pluginConfig
	cfg.Interval = 60 * 1000 - 50
	cfg.Prefix = "u:"

	return cfg, c.GetConfig(&cfg)
}
