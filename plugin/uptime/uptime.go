package uptime

import (
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	//	"gopkg.in/yaml.v1"
	"github.com/op/go-logging"
	"time"
	"os"
	"sync"
	"bufio"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

var log = logging.MustGetLogger("main")

// set up a config struct
type config struct {
	prefix   string
	interval int
}

type plugin struct {
	cfg config
	uptimeFile *os.File
	uptimeFileScanner *bufio.Scanner
	uptimeTpl *util.Template
	uptimeTplShort *util.Template
	dynamicInterval int
	nextTs time.Time
	sync.Mutex
}

type uptimeTpl struct {
	Prefix string
	Uptime string
	UptimeShort string
}

func New(cfg uber.PluginConfig) (uber.Plugin, error) {
	p := &plugin{}
	p.cfg = loadConfig(cfg.Config)
	return  p, nil
}

func (p *plugin) Init() (err error) {
	p.uptimeTpl, err = util.NewTemplate("upt",`{{printf "%s %s" .Prefix .Uptime}}`)
	if err != nil { return }
	p.uptimeTplShort, err = util.NewTemplate("uptimeShort",`u:{{.UptimeShort}}`)
	return
}

func (p *plugin) GetUpdateInterval() int {
	return p.cfg.interval
}
func (p *plugin) UpdatePeriodic() uber.Update {
	var update uber.Update
	uptime := p.getUptime()
	uptimeValue :=  uptimeTpl {
		Prefix: p.cfg.prefix,
		Uptime: util.FormatDuration(uptime),
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
	return p.UpdatePeriodic()
}

// parse received structure into config
func loadConfig(c map[string]interface{}) config {
	var cfg config
	cfg.interval = 60 * 1000 - 50
	cfg.prefix = "ex: "
	for key, value := range c {
		converted, ok := value.(string)
		if ok {
			switch {
			case key == `prefix`:
				cfg.prefix = converted
			default:
				log.Warningf("unknown config key: [%s]", key)

			}
		} else {
			converted, ok := value.(int)
			if ok {
				switch {
				case key == `interval`:
					cfg.interval = converted
				default:
					log.Warningf("unknown config key: [%s]", key)
				}
			} else {
				log.Errorf("Cant interpret value of config key [%s]", key)
			}
		}
	}
	return cfg
}
