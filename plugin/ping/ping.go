package ping

import (
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	"github.com/XANi/golibs/ewma"
	//	"gopkg.in/yaml.v1"
	"github.com/op/go-logging"
	"time"
	"sync"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

var log = logging.MustGetLogger("main")

// set up a config struct
type config struct {
	prefix   string
	interval int
	addrType string
	addr     string
	inflight int
}

type state struct {
	cfg config
	pingCh chan *pingResult
	dropRate *ewma.Ewma
	pingAvg *ewma.Ewma
	stats pingStat
	cnt int
	ev  int
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

func Run(cfg uber.PluginConfig) {
	var st state
	st.cfg = loadConfig(cfg.Config)
	st.dropRate = ewma.NewEwma(time.Duration(15 * time.Second))
	st.pingAvg = ewma.NewEwma(time.Duration(60 * time.Second))
	st.pingCh = make(chan *pingResult,30)
	t := time.Duration(time.Second)
	for i := 0; i < st.cfg.inflight; i++ {
		switch st.cfg.addrType {
		case "tcp":
			go tcpPing(st.cfg.addr, st.pingCh, t)
		case "http":
			go httpPing(st.cfg.addr, st.pingCh, t)
		default:
			log.Panicf("ping: protocol %s not supported", st.cfg.addrType)
		}
	}
	go func() {
		for {
			ping := <- st.pingCh
			st.updateState(ping)
		}
	}()
	// initial update on start
	cfg.Update <- st.updatePeriodic()
	for {
		select {
		// call update when user clicked on the plugin
		case updateEvent := (<-cfg.Events):
			cfg.Update <- st.updateFromEvent(updateEvent)
			// that will wait 10 seconds on no event a
			// and it will "eat" next event to switch to "normal" display
			// basically making it "toggle" between two different views
			select {
			case _ = <-cfg.Events:
				cfg.Update <- st.updatePeriodic()
			case <-time.After(10 * time.Second):
			}
		// update on trigger from main code, this can be used to make all widgets update at the same time if that way is preferred over async
		case _ = <-cfg.Trigger:
			cfg.Update <- st.updatePeriodic()
		// update every interval if nothing triggered update before tat
		case <-time.After(time.Duration(st.cfg.interval) * time.Millisecond):
			cfg.Update <- st.updatePeriodic()
		}
	}
}

func (state *state) updatePeriodic() uber.Update {
	var update uber.Update
	//TODO: cache tpl
	tpl, _ := util.NewTemplate("uberEvent",`{{if not .Ok}}{{color "#aa0000" "png!"}}{{ else }}ping{{end}}: {{formatDuration .LastPing}} {{printf "%2.2f" .DropRate}}%`)
	update.FullText =  tpl.ExecuteString(state.stats)
	update.Markup = `pango`
	update.ShortText = `nope`
	update.Color = util.GetColorPct(int(state.stats.DropRate))
	state.cnt++
	return update
}

func (state *state) updateFromEvent(e uber.Event) uber.Update {
	var update uber.Update

	tpl, _ := util.NewTemplate("uberEvent",`avg: {{formatDuration .AvgPing}}`)
	update.FullText =  tpl.ExecuteString(state.stats)
	update.ShortText = `upd`
	update.Color = `#cccc66`
	state.ev++
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
