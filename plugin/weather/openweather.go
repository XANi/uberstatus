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
	Atmosphere struct{
		Temperature float64 `json:"temp"`
		Pressure float64 `json:"pressure"`
		TemperatureMin float64 `json:"temp_min"`
		TemperatureMax float64 `json:"temp_max"`
	} `json:"main"`
}

func (state *state)  getOpenweatherCurrent () uber.Update {
	var update uber.Update
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	r, err := client.Get(fmt.Sprintf(
		"http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s",
		url.QueryEscape(state.cfg.openWeatherLocation),
		state.cfg.openWeatherApiKey))
	if err != nil {
		log.Warningf("Weather get error: %s", err)
	}
    defer r.Body.Close()
	var weather openweatherCurrentWeather
	err =  json.NewDecoder(r.Body).Decode(&weather)
		if err != nil {
		log.Warningf("Weather JSON decode error: %s", err)
		}
	log.Warningf("out: %+v", weather)
	update.Markup = `pango`
	update.FullText = fmt.Sprintf("T: %2.1f", weather.Atmosphere.Temperature - 273.15)
	return update
}
