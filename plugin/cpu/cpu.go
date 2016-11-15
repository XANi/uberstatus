package cpu

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
	prefix   string
	interval int
	zero     bool
}

type state struct {
	cfg           config
	cnt           int
	ev            int
	previousTicks []cpuTicks
	ticksDiff     []cpuTicks
}

// pregenerate lookup table at start
var ltBar = make(map[int8]string)
var ltColor = make(map[int8]string)

func init() {
	generateLookupTables()
}

func Run(cfg uber.PluginConfig) {
	var st state
	st.cfg = loadConfig(cfg.Config)
	// init with current's total
	st.previousTicks, _ = GetCpuTicks()
	st.ticksDiff, _ = GetCpuTicks()
	// initial update on start
	cfg.Update <- st.updatePeriodic()
	for {
		select {
		case updateEvent := (<-cfg.Events):
			//			cfg.Update <- st.updateFromEvent(updateEvent)
			_ = updateEvent
		case _ = <-cfg.Trigger:
			cfg.Update <- st.updatePeriodic()
		case <-time.After(time.Duration(st.cfg.interval) * time.Millisecond):
			cfg.Update <- st.updatePeriodic()
		}
	}
}

func (state *state) updatePeriodic() uber.Update {
	var update uber.Update
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

	if state.cfg.zero {
		update.FullText = fmt.Sprintf("%s%05.2f%%%s", state.cfg.prefix, usagePct, bars)
	} else {
		update.FullText = fmt.Sprintf("%s% 5.2f%%%s", state.cfg.prefix, usagePct, bars)
	}
	update.ShortText = fmt.Sprintf("%s", util.GetBarChar(int(usagePct)))
	update.Color = util.GetColorPct(int(usagePct))
	update.Markup = `pango`
	state.cnt++
	return update
}

func (state *state) updateFromEvent(e uber.Event) uber.Update {
	var update uber.Update
	update.FullText = fmt.Sprintf("%s %+v", state.cfg.prefix, state.previousTicks)
	update.ShortText = `upd`
	update.Color = `#cccc66`
	state.ev++
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

// parse received structure into config
func loadConfig(c map[string]interface{}) config {
	var cfg config
	cfg.interval = 1000
	cfg.prefix = "ex: "
	cfg.zero = true
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
				converted, ok := value.(bool)
				if ok {
					switch {
					case key == `zero`:
						cfg.zero = converted
					default:
						log.Warningf("unknown config key: [%s]", key)
					}
				} else {
					log.Errorf("Cant interpret value of config key [%s]", key)
				}
			}
		}
	}
	return cfg
}
