package df

import (
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	//	"gopkg.in/yaml.v1"
	"fmt"
	"github.com/op/go-logging"
	"syscall"
	"time"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

var log = logging.MustGetLogger("main")

// set up a config struct
type config struct {
	prefix            string
	interval          int
	mounts            []string
	dfWarningMB       uint64
	dfWarningPercent  uint64
	dfCriticalMB      uint64
	dfCriticalPercent uint64
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
			// that will wait 10 seconds on no event a
			// and it will "eat" next event to switch to "normal" display
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
	update.Markup = `pango`
	update.FullText = fmt.Sprintf(`<span color="#cccccc">%s</span>`, state.cfg.prefix)
	for _, part := range state.cfg.mounts {
		diskFree, diskTotal := getDiskStats(part)
		diskFreePercent := (diskFree * 100) / diskTotal
		diskColor := `#aaffaa`
		if diskFree < state.cfg.dfCriticalMB*1024*1024 || diskFreePercent < state.cfg.dfCriticalPercent {
			diskColor = `#cc3333`
		} else if diskFree < state.cfg.dfWarningMB*1024*1024 || diskFreePercent < state.cfg.dfWarningPercent {
			diskColor = `#cc9966`
		}
		update.FullText = update.FullText + fmt.Sprintf(` <span color="#cccccc">%s:</span><span color="%s">%s</span>`, part, diskColor, util.FormatUnitBytes(int64(diskFree)))
	}
	update.ShortText = `nope`
	update.Color = `#aaaaaa`
	return update
}

func (state *state) updateFromEvent(e uber.Event) uber.Update {
	var update uber.Update
	update.FullText = fmt.Sprintf("%s %+v", state.cfg.prefix, e)
	update.ShortText = `upd`
	update.Color = `#cccc66`
	return update
}

func getDiskStats(path string) (free uint64, total uint64) {
	var stat syscall.Statfs_t
	syscall.Statfs(path, &stat)
	// available blocks * block size
	return stat.Bavail * uint64(stat.Bsize), stat.Blocks * uint64(stat.Bsize)
}

// parse received structure into config
func loadConfig(c map[string]interface{}) config {
	var cfg config
	cfg.interval = 10000
	cfg.prefix = "ex: "
	cfg.dfWarningMB = 2048
	cfg.dfWarningPercent = 7
	cfg.dfCriticalMB = 1024
	cfg.dfCriticalPercent = 5
	for key, value := range c {
		converted, ok := value.(string)
		log.Warningf(`key [%+v] v: [%T]`, converted, value)
		if ok {
			switch {
			case key == `prefix`:
				cfg.prefix = converted
			default:
				log.Warningf("unknown config key: [%s]", key)

			}
		} else if converted, ok := value.([]interface{}); ok {
			switch {
			case key == `mounts`:
				cfg.mounts = make([]string, len(converted))
				for i, v := range converted {
					cfg.mounts[i] = v.(string)
				}
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
