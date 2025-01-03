package cpu

import (
	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	"go.uber.org/zap"

	//	"gopkg.in/yaml.v1"
	"fmt"
	"time"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

// set up a pluginConfig struct
type pluginConfig struct {
	Prefix   string
	Interval int
	Zero     bool
}

type state struct {
	cfg           pluginConfig
	cnt           int
	ev            int
	nextTs        time.Time
	l             *zap.SugaredLogger
	previousTicks []cpuTicks
	ticksDiff     []cpuTicks
}

// pregenerate lookup table at start
var ltColor, ltBar = util.GenerateColorBarLookupTable()

func New(cfg uber.PluginConfig) (z uber.Plugin, err error) {
	p := &state{
		l: cfg.Logger,
	}
	p.cfg, err = loadConfig(cfg.Config)
	return p, err
}

func (st *state) Init() error {
	st.previousTicks, _ = GetCpuTicks()
	st.ticksDiff, _ = GetCpuTicks()
	return nil
}
func (st *state) GetUpdateInterval() int {
	return st.cfg.Interval
}

func (state *state) UpdatePeriodic() uber.Update {
	var update uber.Update
	util.WaitForTs(&state.nextTs)
	currentTicks, _ := GetCpuTicks()
	for i, ticks := range currentTicks {
		state.ticksDiff[i] = ticks.Sub(state.previousTicks[i])
		state.previousTicks[i] = ticks
	}
	usagePct := state.ticksDiff[0].GetCpuUsagePercent()
	bars := ""
	for _, d := range state.ticksDiff[1:] {
		bars = bars + ltBar[int8(d.GetCpuUsagePercent())]
	}

	if state.cfg.Zero {
		update.FullText = fmt.Sprintf("%s%05.2f%%%s", state.cfg.Prefix, usagePct, bars)
	} else {
		update.FullText = fmt.Sprintf("%s%5.2f%%%s", state.cfg.Prefix, usagePct, bars)
	}
	update.ShortText = fmt.Sprintf("%s", util.GetBarChar(int(usagePct)))
	update.Color = util.GetColorPct(int(usagePct))
	update.Markup = `pango`
	state.cnt++
	return update
}

func (state *state) UpdateFromEvent(e uber.Event) uber.Update {
	var update uber.Update
	update.FullText = fmt.Sprintf("%s %+v", state.cfg.Prefix, state.previousTicks)
	update.ShortText = fmt.Sprintf("%s %+v", state.cfg.Prefix, state.previousTicks)
	update.Color = `#cccc66`
	state.ev++
	state.nextTs = time.Now().Add(time.Second * 5)
	return update
}

func generateLookupTables() {
	var i int8
	for i = -1; i <= 101; i++ {
		color := util.GetColorPct(int(i))
		ltColor[i] = color
		ltBar[i] = `<span color="` + color + `">` + util.GetBarChar(int(i)) + `</span>`
	}

}

// parse received structure into pluginConfig
func loadConfig(c config.PluginConfig) (pluginConfig, error) {

	var cfg pluginConfig
	cfg.Interval = 1000
	cfg.Prefix = ""
	cfg.Zero = true
	return cfg, c.GetConfig(&cfg)
}
