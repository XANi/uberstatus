package memory

import (
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	//	"gopkg.in/yaml.v1"
	"fmt"
	"github.com/op/go-logging"
	"time"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

var log = logging.MustGetLogger("main")

// set up a config struct
type config struct {
	interval int
	prefix   string
}

type state struct {
	cfg config
}

func Run(cfg uber.PluginConfig) {
	var st state
	st.cfg = loadConfig(cfg.Config)
	// initial update on start
	cfg.Update <- st.updatePeriodic()
	for {
		select {
		case updateEvent := (<-cfg.Events):
			cfg.Update <- st.updateFromEvent(updateEvent)
			select {
			case _ = <-cfg.Events:
				cfg.Update <- st.updatePeriodic()
			case <-time.After(10 * time.Second):
			}
		case _ = <-cfg.Trigger:
			cfg.Update <- st.updatePeriodic()
		case <-time.After(time.Duration(st.cfg.interval) * time.Millisecond):
			cfg.Update <- st.updatePeriodic()
		}
	}
}

func (state *state) updatePeriodic() uber.Update {
	var update uber.Update
	update.Markup = "pango"
	mem := getMemInfo()
	memFree := mem.Free + mem.Cached + mem.Buffers
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
	swapPct := 100 - ((mem.SwapFree * 100) / mem.SwapTotal)
	update.FullText = fmt.Sprintf(`%s<span color="%s">%s</span><span color="%s">%s</span>`,
		state.cfg.prefix,
		util.GetColorPct(int(swapPct)),
		util.GetBarChar(int(swapPct)),
		util.GetColorPct(int(100-memFreePctForCalc)),
		util.FormatUnitBytes(memFree),
	)
	//		util.GetColorPct(
	update.ShortText = fmt.Sprintf(`<span color="%s">%s%</span>`, util.GetColorPct(int(100-memFreePct)), memFreePct)
	update.Color = `#999999`
	return update
}

func (state *state) updateFromEvent(e uber.Event) uber.Update {
	var update uber.Update
	update.Markup = "pango"
	mem := getMemInfo()
	update.FullText = fmt.Sprintf(`<span color="#bbbbbb">Tot:</span> %s <span color="#bbbbbb">Buf:</span> %s <span color="#bbbbbb">Cache:</span> %s <span color="#bbbbbb">Swap U/C/T</span> %s/%s/%s`,
		util.FormatUnitBytes(mem.Total),
		util.FormatUnitBytes(mem.Buffers),
		util.FormatUnitBytes(mem.Cached),
		util.FormatUnitBytes(mem.SwapTotal-mem.SwapFree),
		util.FormatUnitBytes(mem.SwapCached),
		util.FormatUnitBytes(mem.SwapTotal),
	)
	update.Color = `#999999`
	return update
}

// parse received structure into config
func loadConfig(c map[string]interface{}) config {
	var cfg config
	cfg.interval = 10000
	cfg.prefix = `MEM: `
	for key, value := range c {
		converted, ok := value.(string)
		if ok {
			switch {
			case key == `prefix`:
				cfg.prefix = converted
			default:
				log.Warning("unknown config key: [%s]", key)

			}
		} else {
			converted, ok := value.(int)
			if ok {
				switch {
				case key == `interval`:
					cfg.interval = converted
				default:
					log.Warning("unknown config key: [%s]", key)
				}
			} else {
				log.Error("Cant interpret value of config key [%s]", key)
			}
		}
	}
	return cfg
}
