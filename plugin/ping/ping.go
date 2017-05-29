package ping

import (
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	"github.com/XANi/golibs/ewma"
	//	"gopkg.in/yaml.v1"
	"github.com/op/go-logging"
	"time"
	"sync"
	"fmt"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

var log = logging.MustGetLogger("main")

// set up a config struct
type config struct {
	prefix   string
	interval int
	pingInterval int
	addrType string
	addr     string
	inflight int
}

type state struct {
	cfg config
	dropRate *ewma.Ewma
	pingAvg *ewma.Ewma
	stats *pingStat
	cnt int
	ev  int
	ping func(addr string) pingResult
	nextTs time.Time
	tpl *util.Template
}

type pingResult struct {
	Ok bool
	Duration time.Duration
	OkCount uint64
	FailCount uint64
	DropRate float64
}

type pingStat struct {
	Ok bool
	LastPing time.Duration
	AvgPing time.Duration
	DropRate float64
	sync.Mutex
}

type pinger interface {
	Ping(addr string) *pingResult
}

func New(cfg uber.PluginConfig) (uber.Plugin, error) {
	var st state
	st.cfg = loadConfig(cfg.Config)
	st.dropRate = ewma.NewEwma(time.Duration(15 * time.Second))
	st.pingAvg = ewma.NewEwma(time.Duration(60 * time.Second))
	st.stats = &pingStat{}
	switch st.cfg.addrType {
	case "tcp":
		st.ping = tcpPing
	case "http":
		st.ping = httpPing
	default:
		return &st, fmt.Errorf("ping: protocol %s not supported", st.cfg.addrType)
	}
	var err error
	st.tpl, err = util.NewTemplate("uberEvent",`{{if not .Ok}}{{color "#aa0000" "png!"}}{{ else }}ping{{end}}: {{formatDuration .LastPing}} {{printf "%2.2f" .DropRate}}%`)
	return &st, err
}
func (st *state)Init() error {
	go func() {
		for {
			pingUpd := st.ping(st.cfg.addr)
			st.updateState(&pingUpd)
			if st.cfg.pingInterval > 0 {
				time.Sleep(time.Duration(st.cfg.pingInterval) * time.Millisecond)
			} else {
				time.Sleep(time.Duration(st.cfg.interval) * time.Millisecond)
			}
		}
	}()
	return nil
}
func (state *state) UpdatePeriodic() uber.Update {
	var update uber.Update
	util.WaitForTs(&state.nextTs)
	//TODO: cache tpl
	update.FullText =  state.tpl.ExecuteString(state.stats)
	update.Markup = `pango`
	update.ShortText = `nope`
	update.Color = util.GetColorPct(int(state.stats.DropRate))
	state.cnt++
	return update
}

func (state *state) UpdateFromEvent(e uber.Event) uber.Update {
	var update uber.Update

	tpl, _ := util.NewTemplate("uberEvent",`avg: {{formatDuration .AvgPing}}`)
	update.FullText =  tpl.ExecuteString(state.stats)
	update.ShortText = `upd`
	update.Color = `#cccc66`
	state.ev++
	// display state for at least few seconds before getting to normal update
	state.nextTs = time.Now().Add(time.Second * 3)
	return update
}

func (state *state)updateState(p *pingResult) {
	state.stats.Lock()
	if p.Ok {
		state.stats.DropRate = state.dropRate.UpdateNow(0)
		state.stats.LastPing = p.Duration
		state.stats.AvgPing = time.Duration(int64(state.pingAvg.UpdateNow(float64(p.Duration.Nanoseconds()))))
		state.stats.Ok = p.Ok
	} else  {
		state.stats.DropRate = state.dropRate.UpdateNow(100)
		state.stats.Ok = p.Ok
	}
	state.stats.Unlock()

}
// parse received structure into config
func loadConfig(c map[string]interface{}) config {
	var cfg config
	cfg.interval = 1000
	cfg.addrType = "tcp"
	cfg.addr = "localhost:22"
	cfg.inflight = 10
	for key, value := range c {
		converted, ok := value.(string)
		if ok {
			switch {
			case key == "prefix":
				cfg.prefix = converted
			case key == "type":
				cfg.addrType = converted
			case key == "addr":
				cfg.addr = converted
			default:
				log.Warningf("unknown config key: [%s]", key)

			}
		} else {
			converted, ok := value.(int)
			if ok {
				switch {
				case key == `interval`:
					cfg.interval = converted
				case key == `ping_interval`:
					cfg.pingInterval = converted
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
