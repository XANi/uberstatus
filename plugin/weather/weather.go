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
	openWeatherApiKey string
	openWeatherLocation string
	prefix string
	interval int
}

type state struct {
	cfg config
	cnt int
	ev  int
	currentWeather *openweatherCurrentWeather
	lastWeatherUpdate time.Time
}

type OpenWeatherMapWeather struct {

}

func Run(cfg uber.PluginConfig) {
	var st state
	st.cfg = loadConfig(cfg.Config)
	// initial update on start
	cfg.Update <- st.updatePeriodic()
	for {
		select {
		// call update when user clicked on the plugin
		case updateEvent := (<-cfg.Events):
			cfg.Update <- st.updateFromEvent(updateEvent)
			// that will wait 10 seconds on no event a
			// and it will "eat" next event to switch to "normal" display
			// basically making it "toggle" between two different views
			select {
			case _ = <-cfg.Events:
				cfg.Update <- st.updatePeriodic()
			case <-time.After(10 * time.Second):
			}
		// update on trigger from main code, this can be used to make all widgets update at the same time if that way is preferred over async
		case _ = <-cfg.Trigger:
			cfg.Update <- st.updatePeriodic()
		// update every interval if nothing triggered update before tat
		case <-time.After(time.Duration(st.cfg.interval) * time.Millisecond):
			cfg.Update <- st.updatePeriodic()
		}
	}
}

func (state *state) updatePeriodic() uber.Update {
	return state.getOpenweatherCurrent()
}

func (state *state) updateFromEvent(e uber.Event) uber.Update {
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

func windDirectionToName (deg float64) string {
	if deg > 360 || deg < 0 { return "wind direction outside of 0-360" }
	if deg >= 348.75 || deg < 11.25 { return "N" }
	if deg >= 11.25 && deg < 33.75  { return "NNE"}
	if deg >= 33.75 && deg < 56.25 {return "NE"}
	if deg >= 56.25 && deg < 78.75 { return "ENE"}
	if deg >= 78.75 && deg < 101.25 { return "E"}
	if deg >= 101.25 && deg < 123.75 { return "ESE"}
	if deg >= 123.75 && deg < 146.25 { return "SE"}
	if deg >= 146.25 && deg < 168.75 { return "SSE"}
	if deg >= 168.75 && deg < 191.25 { return "S"}
	if deg >= 191.25 && deg < 213.75 { return "SSW"}
	if deg >= 213.75 && deg < 236.25 {return "SW"}
	if deg >= 236.25 && deg < 258.75 {return "WSW"}
	if deg >= 258.75 && deg < 281.25 {return "W"}
	if deg >= 281.25 && deg < 303.75 {return "WNW"}
	if deg >= 303.75 && deg < 326.25 {return "NW"}
	if deg >= 326.25 && deg < 348.75 {return "NNW"}
	return "ERR: wind direction ranges do not overlap"
}
