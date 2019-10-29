package gpu

import (
	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	//	"gopkg.in/yaml.v1"
	"github.com/op/go-logging"
	"time"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

var log = logging.MustGetLogger("main")

// set up a pluginConfig struct
type pluginConfig struct {
	// some tools allow selecting by pci ID or UUID of device so it is not a number
	GpuID    string
	Prefix   string
	Interval int
}

type plugin struct {
	cfg    pluginConfig
	cnt    int
	ev     int
	nextTs time.Time
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
	var update uber.Update
	// TODO precompile and preallcate
	tpl, _ := util.NewTemplate("uberEvent",`{{color "#00aa00" "Example plugin"}}{{.}}`)
	// example on how to allow UpdateFromEvent to display for some time
	// without being overwritten by periodic updates.
	// We set up ts in our plugin, update it in UpdateFromEvent() and just wait if it is in future via helper function
	util.WaitForTs(&p.nextTs)
	update.FullText =  tpl.ExecuteString(p.cnt)
	update.Markup = `pango`
	update.ShortText = `nope`
	update.Color = `#66cc66`
	p.cnt++
	return update
}

func (p *plugin) UpdateFromEvent(e uber.Event) uber.Update {
	var update uber.Update
	tpl, _ := util.NewTemplate("uberEvent",`{{printf "%+v" .}}`)
	update.FullText =  tpl.ExecuteString(e)
	update.ShortText = `upd`
	update.Color = `#cccc66`
	p.ev++
	// set next TS updatePeriodic will wait to.
	p.nextTs = time.Now().Add(time.Second * 3)
	return update
}
// parse received structure into pluginConfig
func loadConfig(c config.PluginConfig) (pluginConfig,error) {
	var cfg pluginConfig
	cfg.Interval = 60 * 1000 - 50
	cfg.Prefix = "u:"

	return cfg, c.GetConfig(&cfg)
}
