package pipe

import (
	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	//	"gopkg.in/yaml.v1"
	"github.com/op/go-logging"
	"time"
	"os"
	"bufio"
	"sync"
	"fmt"
	"syscall"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

var log = logging.MustGetLogger("main")

// set up a pluginConfig struct
type pluginConfig struct {
	prefix   string
	interval int
	pipePath string
	parseTemplate bool
	markup bool
}

type plugin struct {
	sync.Mutex
	cfg      pluginConfig
	cnt      int
	ev       int
	nextTs   time.Time
	text     string
	updateCh chan uber.Update
}

func New(cfg uber.PluginConfig) (z uber.Plugin,err error) {
	p := &plugin{}
	p.cfg, err = loadConfig(cfg.Config)
	p.updateCh = cfg.Update
	return  p, err
}

func (p *plugin) Init() error {
	go p.startListener()
	return nil
}

func (p *plugin) GetUpdateInterval() int {
	return p.cfg.interval
}
func (p *plugin) UpdatePeriodic() (update uber.Update) {
	if p.cfg.markup {
		update.Markup = `pango`
	} else {
		update.Markup = `none`
	}

	p.Lock()
	defer p.Unlock()
	defer func() {
        if r := recover(); r != nil {
			p.text = fmt.Sprintf("panic in template from pipe [%s]", p.cfg.pipePath )
			update.FullText = p.text
			p.updateCh <- update
        }
	}()
	if p.cfg.parseTemplate {
		tpl, _ := util.NewTemplate("pipe", p.text)
		update.FullText = tpl.ExecuteString(false)
	} else {
		update.FullText = p.text
	}
	if len(update.FullText) > 23 {
		update.ShortText = update.FullText[0:20] + "..."
	} else {
		update.ShortText = update.FullText
	}
	return update
}

func (p *plugin) UpdateFromEvent(e uber.Event) uber.Update {
	return p.UpdatePeriodic()
}

// parse received structure into pluginConfig
func loadConfig(c config.PluginConfig) (pluginConfig,error) {
	var cfg pluginConfig
	cfg.interval = 10000
	cfg.prefix = "ex: "
	cfg.markup = true
	return cfg, c.GetConfig(&cfg)
}

func (p *plugin) startListener() {
	syscall.Mkfifo(p.cfg.pipePath, 0640)
	// pipe needs to be reopened after each writer "disconnects" (EOF)
	for {
		pipe, err := os.OpenFile(p.cfg.pipePath, os.O_RDONLY, 0640)
		if err != nil {
			log.Errorf("Error opening pipe [%s]:%s", p.cfg.pipePath, err)
			return
		}
		scanner := bufio.NewScanner(pipe)
		for scanner.Scan() {
			p.Lock()
			p.text = scanner.Text()
			p.Unlock()
			u := p.UpdatePeriodic()
			p.updateCh <- u
		}
	}

}