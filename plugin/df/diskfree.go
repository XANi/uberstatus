package df

import (
	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	//	"gopkg.in/yaml.v1"
	"fmt"
	"github.com/op/go-logging"
	"syscall"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

var log = logging.MustGetLogger("main")

// set up a config struct
type pluginConfig struct {
	Prefix            string
	Interval          int
	Mounts            []string
	DfWarningMB       uint64
	DfWarningPercent  uint64
	DfCriticalMB      uint64
	DfCriticalPercent uint64
}

type state struct {
	cfg pluginConfig
}
func New(cfg uber.PluginConfig) (z uber.Plugin, err error) {
	p := &state{}
	p.cfg, err = loadConfig(cfg.Config)
	return  p, nil
}
func (st *state)Init() error {
	return nil
}
func (st *state) GetUpdateInterval() int {
	return st.cfg.Interval
}
func (state *state) UpdatePeriodic() uber.Update {
	var update uber.Update
	update.Markup = `pango`
	update.FullText = fmt.Sprintf(`<span color="#cccccc">%s</span>`, state.cfg.Prefix)
	for _, part := range state.cfg.Mounts {
		diskFree, diskTotal := getDiskStats(part)
		if diskTotal == 0 {
			update.FullText = update.FullText + fmt.Sprintf(`<span color="#ffcccc">%s:NaN</span>`,part)
			continue
		}
		diskFreePercent := (diskFree * 100) / diskTotal
		diskColor := `#aaffaa`
		if diskFree < state.cfg.DfCriticalMB*1024*1024 || diskFreePercent < state.cfg.DfCriticalPercent {
			diskColor = `#cc3333`
		} else if diskFree < state.cfg.DfWarningMB*1024*1024 || diskFreePercent < state.cfg.DfWarningPercent {
			diskColor = `#cc9966`
		}
		update.FullText = update.FullText + fmt.Sprintf(` <span color="#cccccc">%s:</span><span color="%s">%s</span>`, part, diskColor, util.FormatUnitBytes(int64(diskFree)))
	}
	update.ShortText = `nope`
	update.Color = `#aaaaaa`
	return update
}

func (state *state) UpdateFromEvent(e uber.Event) uber.Update {
	var update uber.Update
	update.FullText = fmt.Sprintf("%s %+v", state.cfg.Prefix, e)
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
func loadConfig(c config.PluginConfig) (pluginConfig, error) {
	var cfg pluginConfig
	cfg.Interval = 10000
	cfg.Prefix = "ex: "
	cfg.DfWarningMB = 2048
	cfg.DfWarningPercent = 7
	cfg.DfCriticalMB = 1024
	cfg.DfCriticalPercent = 5
	return cfg, c.GetConfig(&cfg)
}
