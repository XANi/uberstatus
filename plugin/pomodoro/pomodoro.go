package pomodoro

import (
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	//	"gopkg.in/yaml.v1"
	"github.com/op/go-logging"
	"time"
	"fmt"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

var log = logging.MustGetLogger("main")

// set up a config struct
type config struct {
	prefix   string
	pomodoroTime int
	shortBreakTime int
	longBreakTime int
}



type plugin struct {
	cfg config
	pomodoroEnd time.Time
	breakEnd time.Time
	nextTs time.Time
	state int
	pomodoros int
}

const (
	stopped = iota
	inPomodor
	inBreakStart
	inShortBreak
	inLongBreak
	inBreakEnd
)

func New(cfg uber.PluginConfig) (uber.Plugin, error) {
	p := &plugin{}
	p.cfg = loadConfig(cfg.Config)
	return  p, nil
}

func (p *plugin) Init() error {
	return nil
}

func (p *plugin) GetUpdateInterval() int {
	return 999
}
func (p *plugin) UpdatePeriodic() uber.Update {
	var update uber.Update
	// TODO precompile and preallcate
	// example on how to allow UpdateFromEvent to display for some time
	// without being overwritten by periodic updates.
	// We set up ts in our plugin, update it in UpdateFromEvent() and just wait if it is in future via helper function

	util.WaitForTs(&p.nextTs)
	update.Markup = `pango`
	update.Color = `#66cc66`
	p.nextStateFromTime()
	switch p.state {
	case stopped:
		p.pomodoroEnd = time.Now().Add(time.Duration(p.cfg.pomodoroTime) * time.Minute)
		update.FullText = "stop, click to start"
		update.Color = `#cccc66`
	case inPomodor:
		diff := p.pomodoroEnd.Sub(time.Now())
		update.FullText = fmt.Sprintf(`<span foreground="#ff0000">🍅</span>: %s`, util.FormatDuration(diff))
		update.Color = `#ccccff`
	case inBreakStart:
		update.FullText = fmt.Sprintf(`<span foreground="#000000" background="#aa0000">%d🍅</span><span background="#aa0000">BREAK:</span>`, p.pomodoros)
	case inBreakEnd:
		update.FullText = `<span foreground="#000000" background="#aa0000">⌛</span><span background="#aa0000">END</span>`
	case inShortBreak:
		diff := p.breakEnd.Sub(time.Now())
		update.FullText = fmt.Sprintf("⌛: %s", util.FormatDuration(diff))
		update.Color = `#ccffcc`
	case inLongBreak:
		diff := p.breakEnd.Sub(time.Now())
		update.FullText = fmt.Sprintf("⏲: %s", util.FormatDuration(diff))
		update.Color = `#ccccff`
	default:
		update.FullText = "wtf"
	}
	return update
}

func (p *plugin) UpdateFromEvent(e uber.Event) uber.Update {
	var update uber.Update
	update.Markup = `pango`
	if e.Button == 1 {
		p.nextStateFromClick()
	}
	switch p.state {
	case stopped:
		update.FullText = "pomodoro stopped"
		update.Color = `#cccc66`
	case inPomodor:
		if time.Now().After(p.pomodoroEnd) {
			update.FullText = "pomodoro ended!"
			update.Color = `#cccc66`
		} else {
			update.FullText = fmt.Sprintf(`<span foreground="#ff0000">🍅</span>: %d`, p.pomodoros)
			update.Color = `#ccccff`
		}
	case inShortBreak:
		update.FullText = "short break start"
	case inLongBreak:
		update.FullText = "long break start"
	default:
		update.FullText = "wtf"
	}
	// set next TS updatePeriodic will wait to.
	p.nextTs = time.Now().Add(time.Second * 3)
	return update
}

func (p *plugin) nextStateFromClick() {
	switch p.state {
	case stopped:
		p.state = inPomodor
		p.pomodoroEnd = time.Now().Add(time.Duration(p.cfg.pomodoroTime) * time.Minute)
	case inPomodor:
	case inBreakStart:
		if (p.pomodoros % 4) == 3 {
			p.breakEnd = time.Now().Add(time.Duration(p.cfg.longBreakTime) * time.Minute)
			p.state = inLongBreak
		} else {
			p.breakEnd = time.Now().Add(time.Duration(p.cfg.shortBreakTime)* time.Minute)
			p.state = inShortBreak
		}
	case inBreakEnd:
		p.pomodoroEnd = time.Now().Add(time.Duration(p.cfg.pomodoroTime) * time.Minute)
		p.state=inPomodor
	case inShortBreak:
	case inLongBreak:
	default:
		log.Warningf("out of state machine: %d", p.state)
		p.state = stopped
	}
}
func (p *plugin) nextStateFromTime() {
	switch p.state {
	case stopped:
	case inPomodor:
		if time.Now().After(p.pomodoroEnd) {
			p.pomodoros++
			p.state = inBreakStart
		}
	case inShortBreak:
		if time.Now().After(p.breakEnd) {
			p.state=inBreakEnd
		}
	case inLongBreak:
		if time.Now().After(p.breakEnd) {
			p.state=inBreakEnd
		}
	}
}
// parse received structure into config
func loadConfig(c map[string]interface{}) config {
	var cfg config
	cfg.prefix = "ex: "
	cfg.pomodoroTime = 25
	cfg.shortBreakTime = 5
	cfg.longBreakTime = 15
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
				case key == `pomodoro_interval`:
					cfg.pomodoroTime = converted
				case key == `short_break_interval`:
					cfg.shortBreakTime = converted
				case key == `long_break_interval`:
					cfg.longBreakTime = converted
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
