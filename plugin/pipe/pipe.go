package pipe

import (
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

// set up a config struct
type config struct {
	prefix   string
	interval int
	pipePath string
	parseTemplate bool
	markup bool
}

type plugin struct {
	sync.Mutex
	cfg config
	cnt int
	ev  int
	nextTs time.Time
	text string
	updateCh chan uber.Update
}

func New(cfg uber.PluginConfig) (uber.Plugin, error) {
	p := &plugin{}
	p.cfg = loadConfig(cfg.Config)
	if len(p.cfg.pipePath) == 0 {
		return p, fmt.Errorf("pipe: path can't be empty")
	}
	p.updateCh = cfg.Update
	return  p, nil
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

// parse received structure into config
func loadConfig(c map[string]interface{}) config {
	var cfg config
	cfg.interval = 10000
	cfg.prefix = "ex: "
	cfg.markup = true
	for key, value := range c {
		switch converted := value.(type) {
		case string:
			switch {
			case key == `path`:
				cfg.pipePath = converted
			default:
				log.Warningf("unknown config key: [%s]", key)

			}
		case int:
				switch {
				case key == `interval`:
					cfg.interval = converted
				default:
					log.Warningf("unknown config key: [%s]", key)
				}
		case bool:
			switch key {
			case "parse_template":
				cfg.parseTemplate = converted
			case "markup":
				cfg.markup = converted
			default:
				log.Warningf("unknown config key: [%s]", key)
			}
		default:
			log.Errorf("Cant interpret value of config key [%s] %+v", key, key)
		}
	}
	return cfg
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