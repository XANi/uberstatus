package weather

import (
	"github.com/XANi/uberstatus/uber"
	"time"
	"net/http"
	"net/url"
	"encoding/json"
	"fmt"
)

// Current weather API
//
//      {
//        "coord": {
//          "lon": -0.13,
//          "lat": 51.51
//        },
//        "weather": [
//          {
//            "id": 300,
//            "main": "Drizzle",
//            "description": "light intensity drizzle",
//            "icon": "09n"
//          },
//          {
//            "id": 500,
//            "main": "Rain",
//            "description": "light rain",
//            "icon": "10n"
//          }
//        ],
//        "base": "stations",
//        "main": {
//          "temp": 278.74,
//          "pressure": 1031,
//          "humidity": 87,
//          "temp_min": 277.15,
//          "temp_max": 280.15
//        },
//        "visibility": 10000,
//        "wind": {
//          "speed": 2.6,
//          "deg": 220
//        },
//        "clouds": {
//          "all": 80
//        },
//        "dt": 1479088200,
//        "sys": {
//          "type": 1,
//          "id": 5091,
//          "message": 0.0235,
//          "country": "GB",
//          "sunrise": 1479107857,
//          "sunset": 1479139909
//        },
//        "id": 2643743,
//        "name": "London",
//        "cod": 200
//      }

type openweatherCurrentWeather struct {
	Weather []struct {
		Id int `json:"id"`
		Name string `json:"name"`
		Description string `json:"description"`
		Icon string `json:"icon"`
	} `json:"weather"`
	Atmosphere openweatherAtmosphere `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
		Direction float64 `json:"deg"`
	}
}

type openweatherAtmosphere struct {
	Temperature float64 `json:"temp"`
	Pressure float64 `json:"pressure"`
	Humidity float64 `json:"humidity"`
	TemperatureMin float64 `json:"temp_min"`
	TemperatureMax float64 `json:"temp_max"`
}

func (state *state) updateWeather () {
	timeout := time.Duration(60 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	r, err := client.Get(fmt.Sprintf(
		"http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s",
		url.QueryEscape(state.cfg.openWeatherLocation),
		state.cfg.openWeatherApiKey))
	if err != nil {
		log.Warningf("Weather get error: %s", err)
		return
	}
    defer r.Body.Close()
	var weather openweatherCurrentWeather
	err =  json.NewDecoder(r.Body).Decode(&weather)
	if err != nil {
		log.Warningf("Weather JSON decode error: %s", err)
	}
	state.currentWeather = &weather
	state.lastWeatherUpdate = time.Now()
}

func (state *state)  getOpenweatherCurrent () uber.Update {
	var update uber.Update
	if !time.Now().Before(state.lastWeatherUpdate.Add(10 * time.Minute)) {
		state.updateWeather()
	}
	update.Markup = `pango`
	// TODO discard old data
	if state.currentWeather != nil {
		update.FullText = parseTemperature(&state.currentWeather.Atmosphere)
		if len ( state.currentWeather.Weather ) > 0 {
			update.FullText = update.FullText + " - " +  state.currentWeather.Weather[0].Description
		}
	} else {
		update.FullText = "cant get weather"
	}
	return update
}

func (state *state)  getOpenweatherPrognosis () uber.Update {
	var update uber.Update
	update.Markup = `pango`
	if state.currentWeather != nil {
		update.FullText = fmt.Sprintf(
			"%1.f hPa, %1.f %%, %1.2f m/s %s",
			state.currentWeather.Atmosphere.Pressure,
			state.currentWeather.Atmosphere.Humidity,
			state.currentWeather.Wind.Speed,
			windDirectionToName(state.currentWeather.Wind.Direction),
		)
	} else {
		update.FullText = "data update failed"
	}
	log.Error(update.FullText)
	return update
}

func parseTemperature(atmosphere *openweatherAtmosphere) string {
	var temperature string
	if (atmosphere.TemperatureMax - atmosphere.TemperatureMin) < 1 {
		temperature = colorizeTemperature(atmosphere.Temperature - 273.15) + "℃"
	} else {
		temperature = fmt.Sprintf(`%s/%s/%s℃`,
			colorizeTemperature(atmosphere.TemperatureMin - 273.15),
			colorizeTemperature(atmosphere.Temperature - 273.15),
			colorizeTemperature(atmosphere.TemperatureMax - 273.15),
		)
	}
	return temperature
}
