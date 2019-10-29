package weather

import (
	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/uber"
	//	"gopkg.in/yaml.v1"
	"github.com/op/go-logging"
	"time"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

var log = logging.MustGetLogger("main")

// set up a pluginConfig struct
type pluginConfig struct {
	OpenWeatherApiKey   string `yaml:"openweather_api_key"`
	OpenWeatherLocation string `yaml:"openweather_location"`
	Prefix              string
	Interval            int
}

type state struct {
	cfg               pluginConfig
	cnt               int
	ev                int
	currentWeather    *openweatherCurrentWeather
	lastWeatherUpdate time.Time

}

type OpenWeatherMapWeather struct {
}

func New(cfg uber.PluginConfig) (z uber.Plugin, err error) {
	p := &state{}
	p.cfg, err = loadConfig(cfg.Config)
	return  p, nil
}

func (state *state) Init() error {
	return nil
}

func (state *state) GetUpdateInterval() int {
	return state.cfg.Interval
}

func (state *state) UpdatePeriodic() uber.Update {
	return state.getOpenweatherCurrent()
}

func (state *state) UpdateFromEvent(e uber.Event) uber.Update {
	return state.getOpenweatherPrognosis()
}

// parse received structure into pluginConfig
func loadConfig(c config.PluginConfig) (pluginConfig,error) {
	var cfg pluginConfig
	cfg.Interval = 60 * 1000 - 50
	cfg.Prefix = "u:"

	return cfg, c.GetConfig(&cfg)
}
