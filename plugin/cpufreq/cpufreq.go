package cpufreq

import (
	"fmt"
	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	"go.uber.org/zap"

	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Cpufreq plugin for uberstatus

var cpufreqRe = regexp.MustCompile(`/sys/devices/system/cpu/cpu(\d+)/cpufreq/scaling_cur_freq`)

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

// pregenerate lookup table at start
var ltColor, ltBar = util.GenerateColorBarLookupTable()

// set up a pluginConfig struct
type pluginConfig struct {
	Prefix   string
	Interval int
}

type plugin struct {
	cfg         pluginConfig
	lowestFreq  int
	highestFreq int
	l           *zap.SugaredLogger
	nextTs      time.Time
}

func New(cfg uber.PluginConfig) (z uber.Plugin, err error) {
	p := &plugin{
		l: cfg.Logger,
	}
	p.lowestFreq = MaxInt
	p.cfg, err = loadConfig(cfg.Config)
	return p, err
}

func (p *plugin) Init() error {
	return nil
}

func (p *plugin) GetUpdateInterval() int {
	return p.cfg.Interval
}
func (p *plugin) UpdatePeriodic() uber.Update {
	var update uber.Update
	// TODO precompile and preallcate
	// cpufreq on how to allow UpdateFromEvent to display for some time
	// without being overwritten by periodic updates.
	// We set up ts in our plugin, update it in UpdateFromEvent() and just wait if it is in future via helper function
	util.WaitForTs(&p.nextTs)
	update.FullText = p.getCpufreqBars()
	update.Markup = `pango`
	update.ShortText = `nope`
	update.Color = `#66cc66`
	return update
}

func (p *plugin) UpdateFromEvent(e uber.Event) uber.Update {
	var update uber.Update
	//	tpl, _ := util.NewTemplate("uberEvent",`{{printf "%+v" .}}`)
	//	tpl.ExecuteString(e)
	update.FullText = fmt.Sprintf("min: %d, max: %d MHz", p.lowestFreq/1000/1000, p.highestFreq/1000/1000)
	update.ShortText = `upd`
	update.Color = `#cccc66`
	// set next TS updatePeriodic will wait to.
	p.nextTs = time.Now().Add(time.Second * 3)
	return update
}

func (p *plugin) getCpufreqBars() string {
	cpufreq, err := p.getCpufreq()
	if err != nil {
		p.l.Errorf("can't get cpufreq: %s", err)
	}
	out := ""
	for _, freq := range cpufreq {
		var pct float32
		zeroOffset := freq - p.lowestFreq
		zeroOffsetMax := p.highestFreq - p.lowestFreq
		if zeroOffset < 0 {
			p.l.Errorf("frequency %d below minfreq %d", freq, p.lowestFreq)
		}
		if zeroOffsetMax > 0 { // if it's zero it means min freq is = max so no scaling. If it is below, wtf
			pct = float32(zeroOffset) / float32(zeroOffsetMax) * 100
		} else {
			pct = 100
		}
		if pct > 100 {
			p.l.Errorf("percent is above 100^ for some reason %d %d %d", freq, p.lowestFreq, p.highestFreq)
		}
		out = out + ltBar[int8(pct)]
	}
	return out
}

// parse received structure into pluginConfig
func loadConfig(c config.PluginConfig) (pluginConfig, error) {
	var cfg pluginConfig
	cfg.Interval = 10000
	cfg.Prefix = "ex: "
	return cfg, c.GetConfig(&cfg)
}

func (p *plugin) getCpufreq() ([]int, error) {
	match, err := filepath.Glob(`/sys/devices/system/cpu/cpu*/cpufreq/scaling_cur_freq`)
	if err != nil {
		return nil, err
	}
	cpufreq := make([]int, len(match))
	for _, file := range match {
		match := cpufreqRe.FindStringSubmatch(file)
		if len(match) > 1 {
			id, _ := strconv.Atoi(match[1])
			content, err := ioutil.ReadFile(file)
			freq, err := strconv.Atoi(strings.TrimSpace(string(content)))
			freq = freq * 1000
			if p.lowestFreq > freq {
				p.lowestFreq = freq
			}
			if p.highestFreq < freq {
				p.highestFreq = freq
			}
			cpufreq[id] = freq
			if err != nil {
				return nil, err
			}
		}
	}
	return cpufreq, nil
}
