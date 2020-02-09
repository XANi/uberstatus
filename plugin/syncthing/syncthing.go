package syncthing

import (
	"sort"
	"sync"

	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	//	"gopkg.in/yaml.v1"
	"github.com/op/go-logging"
	"time"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

var log = logging.MustGetLogger("main")

// set up a pluginConfig struct
type pluginConfig struct {
	Prefix   string
	Interval int
	ApiKey string `yaml:"api_key"`
	ServerAddr string `yaml:"server_addr"`
}

type plugin struct {
	cfg    pluginConfig
	cnt    int
	ev     int
	nextTs time.Time
	folderIdToFolder map[string]string
	folderStatus map[string]FolderStatusId
	folderCompletion map[string]float32
	sync.RWMutex

}




func New(cfg uber.PluginConfig) (z uber.Plugin, err error) {
	p := &plugin{}
	p.cfg, err = loadConfig(cfg.Config)
	p.folderIdToFolder = make(map[string]string)
	p.folderStatus = make(map[string]FolderStatusId)
	p.folderCompletion = make(map[string]float32)
	if len(p.cfg.ServerAddr) == 0 {
		p.cfg.ServerAddr = "http://127.0.0.1:8384"
	}
	return  p, err
}

func (p *plugin) Init() error {
	go func () {
		for {
			// most of our updates will come from events
			// on top of that syncthing folder state is blocking API...
			time.Sleep(time.Second * 120)
			p.updateSyncthingFolders()
		}
	} ()
	go func() {
		p.updateSyncthingFolders()
		p.updateSynctingEvents()
	} ()

	return nil
}

func (p *plugin) GetUpdateInterval() int {
	return p.cfg.Interval
}
func (p *plugin) UpdatePeriodic() uber.Update {
	var update uber.Update

	out := p.cfg.Prefix

	p.RLock()
	var list  []string
	for id, _ := range p.folderIdToFolder {
		list = append(list,id)
	}
	sort.Slice(list, func(i, j int) bool { return p.folderIdToFolder[list[i]] < p.folderIdToFolder[list[j]] })
	for _, id := range list {
		if state, ok := p.folderStatus[id];ok  {
			switch state {
			case StatusIdle: out += `{{color "#aaaaff" "█"}}`
			case StatusScanning: out += `{{color "#aaffaa" "C"}}`
			case StatusScanWaiting: out += `{{color "#aaffff" "W"}}`
			case StatusSyncing:
				out += `{{color "#ffffaa" "`
				out += util.GetBarChar(int(p.folderCompletion[id]))
				out += `"}}`
			case StatusSynPreparing: out += `{{color "#ffff00" "P"}}`
			case StatusUnknown: out += `{{color "#aaaaaa" "#"}}`
			default: out += `{{color "#aaaaaa" "?"}}`
			}
		} else {
			out += `{{color "#ffaaaa" "?"}}`
		}
	}
	p.RUnlock()
	tpl, _ := util.NewTemplate("uberEvent",out)
	// example on how to allow UpdateFromEvent to display for some time
	// without being overwritten by periodic updates.
	// We set up ts in our plugin, update it in UpdateFromEvent() and just wait if it is in future via helper function
	util.WaitForTs(&p.nextTs)
	update.FullText =  tpl.ExecuteString(p.folderStatus)
	update.Markup = `pango`
	update.ShortText = `nope`
	update.Color = `#66cc66`
	p.cnt++

	return update
}

func (p *plugin) UpdateFromEvent(e uber.Event) uber.Update {
	var update uber.Update
	tpl, _ := util.NewTemplate("uberEvent",`{{printf "%+v" .}}`)
	update.FullText =  tpl.ExecuteString(p.folderStatus)
	update.ShortText = `upd`
	update.Color = `#cccc66`
	p.ev++
	// set next TS updatePeriodic will wait to.
	p.nextTs = time.Now().Add(time.Second * 3)
	return update
}

// parse received structure into pluginConfig
func loadConfig(c config.PluginConfig) (pluginConfig,error) {
	var cfg pluginConfig
	cfg.Interval = 1000
	cfg.Prefix = "S: "
	// optionally, check for pluginConfig validity after GetConfig call
	return cfg, c.GetConfig(&cfg)
}
