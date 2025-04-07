package icinga

import (
	"fmt"
	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	"github.com/efigence/go-icinga2"
	"github.com/efigence/go-mon"
	"go.uber.org/zap"
	"sync"

	"time"
)

// Icinga2 monitoring

type pluginConfig struct {
	Prefix            string
	Interval          int
	URL               string        `yaml:"url"`
	User              string        `yaml:"user"`
	Pass              string        `yaml:"pass"`
	HostFilter        string        `yaml:"host_filter"`
	ServiceFilter     string        `yaml:"service_filter"`
	APIUpdateInterval time.Duration `yaml:"api_update_interval"`
}

type plugin struct {
	l                  *zap.SugaredLogger
	i                  *icinga2.API
	cfg                pluginConfig
	hostStatusMap      map[string]int
	serviceStatusMap   map[string]int
	serviceHardnessMap map[string]int
	servicesDown       []string
	lastErr            error
	nextTs             time.Time
	sync.Mutex
}

func New(cfg uber.PluginConfig) (z uber.Plugin, err error) {
	p := &plugin{
		l: cfg.Logger,
	}
	p.cfg, err = loadConfig(cfg.Config)
	return p, err
}

func (p *plugin) Init() (err error) {
	p.i, err = icinga2.New(p.cfg.URL, p.cfg.User, p.cfg.Pass)
	if p.cfg.APIUpdateInterval < time.Second {
		p.cfg.APIUpdateInterval = time.Second * 30
	}
	p.hostStatusMap = map[string]int{
		"invalid":     0,
		"up":          0,
		"down":        0,
		"unreachable": 0,
	}
	p.serviceStatusMap = map[string]int{
		"invalid":  0,
		"ok":       0,
		"warning":  0,
		"critical": 0,
		"unknown":  0,
	}

	p.servicesDown = []string{}
	w := make(chan bool)
	go func() {
		p.update()
		w <- true
		for {
			time.Sleep(p.cfg.APIUpdateInterval)
			p.update()
		}
	}()
	// wait for a sec so we might get status right away instead of starting plugin with no state.
	select {
	case <-time.After(time.Second):
	case <-w:
	}
	return err
}

func (p *plugin) GetUpdateInterval() int {
	return p.cfg.Interval
}
func (p *plugin) UpdatePeriodic() uber.Update {
	var update uber.Update
	tpl, _ := util.NewTemplate("uberEvent", `{{printf "%+v" .}}`)
	// example on how to allow UpdateFromEvent to display for some time
	// without being overwritten by periodic updates.
	// We set up ts in our plugin, update it in UpdateFromEvent() and just wait if it is in future via helper function
	util.WaitForTs(&p.nextTs)
	p.Lock()
	defer p.Unlock()
	update.Markup = `pango`
	if p.lastErr != nil {
		tpl, _ = util.NewTemplate("uberEvent", `err {{.}}`)
		update.FullText = tpl.ExecuteString(p.lastErr)
		update.ShortText = fmt.Sprintf("%s", p.lastErr)
		update.Color = `#ffcc66`
	} else {
		tpl, err := util.NewTemplate("uberEvent", `OK:{{ index . "ok"}} W:{{ index . "warning"}} C:{{ index . "critical" }}`)
		if err != nil {
			panic(err)
		}
		update.FullText = tpl.ExecuteString(p.serviceStatusMap)
		update.ShortText = tpl.ExecuteString(p.serviceStatusMap)
		if p.serviceHardnessMap["critical"] > 0 {
			update.Color = `#ff6666`
		} else if p.serviceHardnessMap["warning"] > 0 {
			update.Color = `#cccc66`
		} else if p.serviceHardnessMap["unknown"] > 0 {
			update.Color = `#6666cc`
		} else {
			update.Color = `#66cc66`
		}
	}
	return update
}

func (p *plugin) update() {
	hosts, err := p.i.GetHostsByFilter(p.cfg.HostFilter)
	statusCtr := map[uint8]int{}
	for _, h := range hosts {
		statusCtr[h.State]++
	}
	statusMap := map[string]int{
		"invalid":     statusCtr[uint8(mon.HostInvalid)],
		"up":          statusCtr[uint8(mon.HostUp)],
		"down":        statusCtr[uint8(mon.HostDown)],
		"unreachable": statusCtr[uint8(mon.HostUnreachable)],
	}
	p.Lock()
	p.hostStatusMap = statusMap
	if err != nil {
		p.lastErr = err
	}
	p.Unlock()

	services, err2 := p.i.GetServicesByFilter(p.cfg.ServiceFilter)
	statusCtr = map[uint8]int{}
	hardCtr := map[uint8]int{}
	servicesNotOk := []string{}
	for _, s := range services {
		if s.StateHard {
			hardCtr[s.State]++
		}
		statusCtr[s.State]++
		if s.State != uint8(mon.StateOk) {
			servicesNotOk = append(servicesNotOk, fmt.Sprintf("%s:%s", s.Host, s.Service))
		}

	}
	statusMap = map[string]int{
		"invalid":  statusCtr[uint8(mon.StateInvalid)],
		"ok":       statusCtr[uint8(mon.StateOk)],
		"warning":  statusCtr[uint8(mon.StateWarning)],
		"critical": statusCtr[uint8(mon.StateCritical)],
		"unknown":  statusCtr[uint8(mon.StateUnknown)],
	}
	hardMap := map[string]int{
		"invalid":  hardCtr[uint8(mon.StateInvalid)],
		"ok":       hardCtr[uint8(mon.StateOk)],
		"warning":  hardCtr[uint8(mon.StateWarning)],
		"critical": hardCtr[uint8(mon.StateCritical)],
		"unknown":  hardCtr[uint8(mon.StateUnknown)],
	}
	p.Lock()
	p.serviceStatusMap = statusMap
	p.serviceHardnessMap = hardMap
	p.servicesDown = servicesNotOk
	if err != nil {
		p.lastErr = err
	} else {
		p.lastErr = err2
	}
	p.Unlock()
}

func (p *plugin) UpdateFromEvent(e uber.Event) uber.Update {
	var update uber.Update
	tpl, _ := util.NewTemplate("uberEvent", `{{printf "%+v" .}}`)
	// example on how to allow UpdateFromEvent to display for some time
	// without being overwritten by periodic updates.
	// We set up ts in our plugin, update it in UpdateFromEvent() and just wait if it is in future via helper function
	util.WaitForTs(&p.nextTs)
	p.Lock()
	defer p.Unlock()
	update.Markup = `pango`
	if p.lastErr != nil {
		tpl, _ = util.NewTemplate("uberEvent", `err {{.}}`)
		update.FullText = tpl.ExecuteString(p.lastErr)
		update.ShortText = fmt.Sprintf("%s", p.lastErr)
		update.Color = `#ffcc66`
	} else {
		tpl, err := util.NewTemplate("uberEvent", `{{ printf "%+v" . }}`)
		if err != nil {
			panic(err)
		}
		update.FullText = tpl.ExecuteString(p.servicesDown)
		update.ShortText = tpl.ExecuteString(p.servicesDown)
		if p.serviceStatusMap["critical"] > 0 {
			update.Color = `#ff6666`
		} else if p.serviceStatusMap["warning"] > 0 {
			update.Color = `#cccc66`
		} else if p.serviceStatusMap["unknown"] > 0 {
			update.Color = `#6666cc`
		} else {
			update.Color = `#66cc66`
		}
	}

	// set next TS updatePeriodic will wait to.
	p.nextTs = time.Now().Add(time.Second * 3)
	return update
}

// parse received structure into pluginConfig
func loadConfig(c config.PluginConfig) (pluginConfig, error) {
	var cfg pluginConfig
	cfg.Interval = 10000
	cfg.Prefix = "ex: "
	// optionally, check for pluginConfig validity after GetConfig call
	return cfg, c.GetConfig(&cfg)
}
