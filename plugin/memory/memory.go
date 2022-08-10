package memory

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

// set up a config struct
type pluginConfig struct {
	Interval int
	Prefix   string
}

type state struct {
	l               *zap.SugaredLogger
	cfg             pluginConfig
	nextTs          time.Time
	hasMemAvailable bool //only newer kernels have it
}

func New(cfg uber.PluginConfig) (u uber.Plugin, err error) {
	s := &state{
		l: cfg.Logger,
	}
	s.cfg, err = loadConfig(cfg.Config)
	return s, err
}
func (state *state) Init() error {
	mem := getMemInfo()
	if mem.HasAvailable {
		state.l.Info(`has MemAvailable in /proc/meminfo, using that as source for "free" memory`)
	}
	return nil
}

func (state *state) GetUpdateInterval() int {
	return state.cfg.Interval
}

func (state *state) UpdatePeriodic() uber.Update {
	util.WaitForTs(&state.nextTs)
	var update uber.Update
	update.Markup = "pango"
	mem := getMemInfo()
	var memFree int64
	if mem.HasAvailable {
		memFree = mem.Available
	} else {
		memFree = mem.Free + mem.Cached + mem.Buffers
	}
	// some adjustments for high/low mem systems
	// rescale % scale based on total memory
	var memFreePctForCalc float64
	memFreePct := float64(memFree) / float64(mem.Total) * 100
	// cap "total for percent calculation" on 8G
	memTotalForCalc := int64(8192 * 1024 * 1024) // fake total used for free % calculation
	// cap "total for percent calculation" on 8G
	if mem.Total < memTotalForCalc {
		memTotalForCalc = mem.Total
	}
	if memFree > memTotalForCalc {
		memFreePctForCalc = 100
	} else {
		memFreePctForCalc = float64(memFree) / float64(memTotalForCalc) * 100
	}
	var swapPct int64
	if mem.SwapTotal == 0 {
		swapPct = 0
	} else {
		swapPct = 100 - ((mem.SwapFree * 100) / mem.SwapTotal)
	}
	update.FullText = fmt.Sprintf(`%s<span color="%s">%s</span><span color="%s">%s</span>`,
		state.cfg.Prefix,
		util.GetColorPct(int(swapPct)),
		util.GetBarChar(int(swapPct)),
		util.GetColorPct(int(100-memFreePctForCalc)),
		util.FormatUnitBytes(memFree),
	)
	//		util.GetColorPct(
	update.ShortText = fmt.Sprintf(`<span color="%s">%2.f%%</span>`, util.GetColorPct(int(100-memFreePct)), memFreePct)
	update.Color = `#999999`
	return update
}

func (state *state) UpdateFromEvent(e uber.Event) uber.Update {
	var update uber.Update
	update.Markup = "pango"
	mem := getMemInfo()
	update.FullText = fmt.Sprintf(`<span color="#bbbbbb">Tot:</span> %s <span color="#bbbbbb">Buf:</span> %s <span color="#bbbbbb">Cache:</span> %s`,
		util.FormatUnitBytes(mem.Total),
		util.FormatUnitBytes(mem.Buffers),
		util.FormatUnitBytes(mem.Cached),
	)
	if mem.SwapTotal > 0 {
		update.FullText = update.FullText + fmt.Sprintf(` <span color="#bbbbbb">Swap U/C/T</span> %s/%s/%s`,
			util.FormatUnitBytes(mem.SwapTotal-mem.SwapFree),
			util.FormatUnitBytes(mem.SwapCached),
			util.FormatUnitBytes(mem.SwapTotal),
		)
	} else {
		update.FullText = update.FullText + ` <span color="#bb0000">Swap off</span>`
	}
	update.ShortText = update.FullText
	update.Color = `#999999`
	state.nextTs = time.Now().Add(time.Second * 3)
	return update
}

// parse received structure into config
func loadConfig(c config.PluginConfig) (pluginConfig, error) {
	var cfg pluginConfig
	cfg.Interval = 10000
	cfg.Prefix = `MEM: `
	return cfg, c.GetConfig(&cfg)
}
