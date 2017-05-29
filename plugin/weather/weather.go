package weather

import (
	"github.com/XANi/uberstatus/uber"
	//	"gopkg.in/yaml.v1"
	"github.com/op/go-logging"
	"time"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

var log = logging.MustGetLogger("main")

// set up a config struct
type config struct {
	openWeatherApiKey   string
	openWeatherLocation string
	prefix              string
	interval            int
}

type state struct {
	cfg               config
	cnt               int
	ev                int
	currentWeather    *openweatherCurrentWeather
	lastWeatherUpdate time.Time

}

type OpenWeatherMapWeather struct {
}

func New(cfg uber.PluginConfig) (uber.Plugin, error) {
	var st state
	st.cfg = loadConfig(cfg.Config)
	return &st, nil
}

func (state *state) Init() error {
	return nil
}

func (state *state) GetUpdateInterval() int {
	return state.cfg.interval
}

func (state *state) UpdatePeriodic() uber.Update {
	return state.getOpenweatherCurrent()
}

func (state *state) UpdateFromEvent(e uber.Event) uber.Update {
	return state.getOpenweatherPrognosis()
}

// parse received structure into config
func loadConfig(c map[string]interface{}) config {
	var cfg config
	cfg.interval = 1000 * 10 * 1
	cfg.prefix = "ex: "
	for key, value := range c {
		converted, ok := value.(string)
		if ok {
			switch {
			case key == `prefix`:
				cfg.prefix = converted
			case key == `openweather_api_key`:
				cfg.openWeatherApiKey = converted
			case key == `openweather_location`:
				cfg.openWeatherLocation = converted
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
