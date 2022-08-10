package ping

import (
	"fmt"
	"github.com/XANi/golibs/ewma"
	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	"go.uber.org/zap"

	"sync"
	"time"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

// set up a pluginConfig struct
type pluginConfig struct {
	Prefix       string
	Interval     int
	PingInterval int
	AddrType     string
	Addr         string
	Inflight     int
}

type state struct {
	l              *zap.SugaredLogger
	cfg            pluginConfig
	dropRate       *ewma.Ewma
	pingAvg        *ewma.Ewma
	rollingPingAvg *ewma.Ewma
	stats          *pingStat
	cnt            int
	ev             int
	ping           func(addr string) pingResult
	nextTs         time.Time
	tpl            *util.Template
}

type pingResult struct {
	Ok        bool
	Duration  time.Duration
	OkCount   uint64
	FailCount uint64
	DropRate  float64
}

type pingStat struct {
	Ok       bool
	LastPing time.Duration
	AvgPing  time.Duration
	// Drop rate in 0-100%
	DropRate float64
	sync.Mutex
}

type pinger interface {
	Ping(addr string) *pingResult
}

func New(cfg uber.PluginConfig) (z uber.Plugin, err error) {
	var st state
	st.cfg, err = loadConfig(cfg.Config)
	if err != nil {
		return nil, err
	}
	st.dropRate = ewma.NewEwma(time.Duration(15 * time.Second))
	st.pingAvg = ewma.NewEwma(time.Duration(60 * time.Second))
	st.rollingPingAvg = ewma.NewEwma(time.Duration(5 * time.Second))
	st.stats = &pingStat{}
	st.l = cfg.Logger
	switch st.cfg.AddrType {
	case "tcp":
		st.ping = tcpPing
	case "http":
		st.ping = httpPing
	default:
		return &st, fmt.Errorf("ping: protocol %s not supported", st.cfg.AddrType)
	}
	st.tpl, err = util.NewTemplate("uberEvent", `{{if not .Ok}}{{color "#aa0000" "png!"}}{{ else }}ping{{end}}: {{formatDurationPadded .LastPing}} {{printf "%2.2f" .DropRate}}%`)
	return &st, err
}
func (st *state) Init() error {
	go func() {
		for {
			pingUpdCh := make(chan bool, 1)
			go func() {
				upd := st.ping(st.cfg.Addr)
				st.updateState(&upd)
				st.rollingPingAvg.UpdateNow(float64(upd.Duration.Nanoseconds()))
				pingUpdCh <- true
			}()
			done := false
			ts := time.Now()
			i := 0
			for {
				i++
				select {
				case <-time.After(time.Millisecond * 100):
					st.rollingPingAvg.UpdateNow(float64(time.Since(ts).Nanoseconds()))
				case <-pingUpdCh:
					done = true
				}
				if done || i > 20 {
					break
				}
			}

			if st.cfg.PingInterval > 0 {
				time.Sleep(time.Duration(st.cfg.PingInterval) * time.Millisecond)
			} else {
				time.Sleep(time.Duration(st.cfg.Interval) * time.Millisecond)
			}
		}
	}()
	return nil
}
func (st *state) GetUpdateInterval() int {
	return st.cfg.Interval
}
func (state *state) UpdatePeriodic() uber.Update {
	var update uber.Update
	util.WaitForTs(&state.nextTs)
	//TODO: cache tpl
	update.FullText = state.tpl.ExecuteString(state.stats)
	update.Markup = `pango`
	update.ShortText = `nope`
	update.Color = util.GetColorPct(state.CalculateOkayness())
	state.cnt++
	return update
}

func (state *state) CalculateOkayness() int {
	avg := float64(state.stats.AvgPing.Nanoseconds())
	var ratio float64
	var pct float64
	if state.rollingPingAvg.Current > avg*2 {
		ratio = state.rollingPingAvg.Current / (avg + 1)
		pct := 20 * ratio
		if pct > 100 {
			pct = 100
		}
	}
	if state.stats.DropRate > pct {
		return int(state.stats.DropRate)
	} else {
		return int(pct)
	}
}

func (state *state) UpdateFromEvent(e uber.Event) uber.Update {
	var update uber.Update

	tpl, _ := util.NewTemplate("uberEvent", `avg: {{formatDuration .AvgPing}}`)
	update.FullText = tpl.ExecuteString(state.stats)
	update.ShortText = `upd`
	update.Color = `#cccc66`
	state.ev++
	// display state for at least few seconds before getting to normal update
	state.nextTs = time.Now().Add(time.Second * 3)
	return update
}

func (state *state) updateState(p *pingResult) {
	state.stats.Lock()
	if p.Ok {
		state.stats.DropRate = state.dropRate.UpdateNow(0)
		state.stats.LastPing = p.Duration
		state.stats.AvgPing = time.Duration(int64(state.pingAvg.UpdateNow(float64(p.Duration.Nanoseconds()))))
		state.stats.Ok = p.Ok
	} else {
		state.stats.DropRate = state.dropRate.UpdateNow(100)
		state.stats.Ok = p.Ok
	}
	state.stats.Unlock()

}

// parse received structure into pluginConfig
func loadConfig(c config.PluginConfig) (pluginConfig, error) {
	var cfg pluginConfig
	cfg.Interval = 1000
	cfg.AddrType = "tcp"
	cfg.Addr = "localhost:22"
	cfg.Inflight = 10
	return cfg, c.GetConfig(&cfg)
}
