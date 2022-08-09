package gpu

import (
	"bufio"
	"fmt"
	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	//	"gopkg.in/yaml.v1"
	"github.com/op/go-logging"
	"time"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

var log = logging.MustGetLogger("main")

// set up a pluginConfig struct
type pluginConfig struct {
	// some tools allow selecting by pci ID or UUID of device so it is not a number
	ID       string
	Type     string
	Prefix   string
	Interval int
}

// all clocks in MHz
type gpuInfo struct {
	Power         float32
	ClockGraphics int
	// StreamingMultiprocessor
	ClockSM     int
	ClockMemory int
	ClockVideo  int
}

type plugin struct {
	cfg        pluginConfig
	cnt        int
	ev         int
	gpuMax     gpuInfo
	gpuCurrent gpuInfo
	sync.Mutex
	nextTs time.Time
}
type csvColumn struct {
	SmiName string
	// no point doing reflections with 2 values, just choose between.
	// convert to int and put under this
	IntPtr *int
	// convert to float and put under this
	FloatPtr *float32
}

func New(cfg uber.PluginConfig) (z uber.Plugin, err error) {
	p := &plugin{}
	p.cfg, err = loadConfig(cfg.Config)
	if p.cfg.Type != "nvidia" {
		return nil, fmt.Errorf(`only "nvidia" gpu type is supported`)
	}
	return p, nil
}

func (p *plugin) Init() (err error) {
	p.gpuMax, err = nvGetMaximumValues(p.cfg.ID)
	// TODO make GPU specific functions
	go func() {
		// nvidia-smi --help-query-gpu # for list
		var csvColumns = []csvColumn{
			{
				SmiName:  "power.draw",
				FloatPtr: &p.gpuCurrent.Power,
			},
			{
				SmiName:  "enforced.power.limit",
				FloatPtr: &p.gpuMax.Power,
			},
			{
				SmiName: "clocks.current.graphics",
				IntPtr:  &p.gpuCurrent.ClockGraphics,
			},
			{
				SmiName: "clocks.max.graphics",
				IntPtr:  &p.gpuMax.ClockGraphics,
			},
			{
				SmiName: "clocks.current.sm",
				IntPtr:  &p.gpuCurrent.ClockSM,
			},
			{
				SmiName: "clocks.max.sm",
				IntPtr:  &p.gpuMax.ClockSM,
			},
			{
				SmiName: "clocks.current.memory",
				IntPtr:  &p.gpuCurrent.ClockMemory,
			},
			{
				SmiName: "clocks.max.memory",
				IntPtr:  &p.gpuMax.ClockMemory,
			},
			// max video is not available for query, we get it from XML somewhere else
			{
				SmiName: "clocks.current.video",
				IntPtr:  &p.gpuCurrent.ClockVideo,
			},
		}
		var columns []string
		for _, c := range csvColumns {
			columns = append(columns, c.SmiName)
		}
		cmd := exec.Command("nvidia-smi",
			"--query-gpu="+strings.Join(columns, ","),
			"--format=csv,noheader,nounits",
			"--id", p.cfg.ID,
			"-l", "1",
		)
		pipe, err := cmd.StdoutPipe()
		if err != nil {
			log.Errorf("could not connect input pipe to nvidia-smi:%s", err)
		}
		err = cmd.Start()
		if err != nil {
			log.Errorf("could not start nvidia-smi:%s", err)
		}

		scanner := bufio.NewScanner(pipe)
		for scanner.Scan() {
			//log.Errorf("--- %s", scanner.Text())
			ssv := strings.Split(scanner.Text(), ",")
			p.Lock()
			for idx, column := range csvColumns {
				val := strings.TrimSpace(ssv[idx])
				if column.IntPtr != nil {
					i, err := strconv.Atoi(val)
					if err != nil {
						log.Warningf("Error converting [%s] to int; %s", val, err)
					} else {
						*column.IntPtr = i
					}
				}
				if column.FloatPtr != nil {
					f, err := strconv.ParseFloat(val, 32)
					if err != nil {
						log.Warningf("Error converting [%s] to float; %s", val, err)
					} else {
						*column.FloatPtr = float32(f)
					}
				}
			}
			p.Unlock()
			_ = ssv
		}
	}()
	return err
}

func (p *plugin) GetUpdateInterval() int {
	return p.cfg.Interval
}
func (p *plugin) UpdatePeriodic() uber.Update {
	var update uber.Update
	// TODO precompile and preallcate
	tpl, _ := util.NewTemplate("uberEvent", `{{.Power}} W 🕐: {{.ClockGraphics}}/{{.ClockSM}}/{{.ClockMemory}}/{{.ClockVideo}}`)
	// example on how to allow UpdateFromEvent to display for some time
	// without being overwritten by periodic updates.
	// We set up ts in our plugin, update it in UpdateFromEvent() and just wait if it is in future via helper function
	util.WaitForTs(&p.nextTs)
	p.Lock()
	update.FullText = tpl.ExecuteString(p.gpuCurrent) +
		" " +
		util.GetBarChar(int(100*(float32(p.gpuCurrent.ClockGraphics)/float32(p.gpuMax.ClockGraphics+1)))) +
		util.GetBarChar(int(100*(float32(p.gpuCurrent.ClockMemory)/float32(p.gpuMax.ClockMemory+1))))
	p.Unlock()
	update.Markup = `pango`
	update.ShortText = `nope`
	update.Color = `#66cc66`
	p.cnt++
	return update
}

func (p *plugin) UpdateFromEvent(e uber.Event) uber.Update {
	var update uber.Update
	tpl, _ := util.NewTemplate("uberEvent", `MAX: {{.Power}} 🕐: {{.ClockGraphics}}/{{.ClockSM}}/{{.ClockMemory}}/{{.ClockVideo}}`)
	update.FullText = tpl.ExecuteString(p.gpuMax)
	update.ShortText = `upd`
	update.Color = `#cccc66`
	p.ev++
	// set next TS updatePeriodic will wait to.
	p.nextTs = time.Now().Add(time.Second * 3)
	return update
}

// parse received structure into pluginConfig
func loadConfig(c config.PluginConfig) (pluginConfig, error) {
	var cfg pluginConfig
	cfg.Interval = 1000
	cfg.Prefix = "u:"
	err := c.GetConfig(&cfg)
	cfg.Type = strings.ToLower(cfg.Type)
	return cfg, err
}
