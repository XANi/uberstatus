package cpu

import (
	"github.com/XANi/uberstatus/uber"
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
}

type state struct {
	cfg                config
	cnt                int
	ev                 int
	previousTicks      []cpuTicks
	ticksDiff          []cpuTicks
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
		state.previousTicks[i]=ticks
	}
	usagePct := state.ticksDiff[0].GetCpuUsagePercent()
	bars := ""
	for _, d := range state.ticksDiff[1:] {
		bars = bars + ltBar[int8(d.GetCpuUsagePercent())]
	}

	update.FullText = fmt.Sprintf("%05.2f%%%s", usagePct, bars)
	update.ShortText = fmt.Sprintf("%s", getBarChar(usagePct))
	update.Color = getColor(usagePct)
	update.Markup = `pango`;
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



func getBarChar(pct float64) string {
	switch {
	case pct > 90:
		return `█`
	case pct > 80:
		return `▇`
	case pct > 70:
		return `▆`
	case pct > 60:
		return `▅`
	case pct > 40:
		return `▄`
	case pct > 20:
		return `▂`
	case pct > 10:
		return `▁`
	}
	return ` `
}

func getColor(pct float64) string {
	switch {
	case pct > 90:
		return `#dd0000`
	case pct > 80:
		return `#cc3333`
	case pct > 70:
		return `#ccaa44`
	case pct > 50:
		return `#cc9966`
	case pct > 30:
		return `#cccc66`
	case pct > 15:
		return `#66cc66`
	}
	return `#666666`
}

//var lookupTableBar = generateBarLookupTable()
//var lookupTableColor = generateColorLookupTable()

func generateLookupTables()  {
	var i int8
	for i = -1; i <= 101; i++ {
		color := getColor(float64(i))
		ltColor[i] = color
		ltBar[i] = `<span color="` + color + `">` + getBarChar(float64(i)) + `</span>`
	}

}

// parse received structure into config
func loadConfig(c map[string]interface{}) config {
	var cfg config
	cfg.interval = 1000
	cfg.prefix = "ex: "
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
