package weather

import (
	"fmt"
)

func windDirectionToName(deg float64) string {
	if deg > 360 || deg < 0 {
		return "wind direction outside of 0-360"
	}
	if deg >= 348.75 || deg < 11.25 {
		return "N"
	}
	if deg >= 11.25 && deg < 33.75 {
		return "NNE"
	}
	if deg >= 33.75 && deg < 56.25 {
		return "NE"
	}
	if deg >= 56.25 && deg < 78.75 {
		return "ENE"
	}
	if deg >= 78.75 && deg < 101.25 {
		return "E"
	}
	if deg >= 101.25 && deg < 123.75 {
		return "ESE"
	}
	if deg >= 123.75 && deg < 146.25 {
		return "SE"
	}
	if deg >= 146.25 && deg < 168.75 {
		return "SSE"
	}
	if deg >= 168.75 && deg < 191.25 {
		return "S"
	}
	if deg >= 191.25 && deg < 213.75 {
		return "SSW"
	}
	if deg >= 213.75 && deg < 236.25 {
		return "SW"
	}
	if deg >= 236.25 && deg < 258.75 {
		return "WSW"
	}
	if deg >= 258.75 && deg < 281.25 {
		return "W"
	}
	if deg >= 281.25 && deg < 303.75 {
		return "WNW"
	}
	if deg >= 303.75 && deg < 326.25 {
		return "NW"
	}
	if deg >= 326.25 && deg < 348.75 {
		return "NNW"
	}
	return "ERR: wind direction ranges do not overlap"
}

func windDirectionToArrow(deg float64) string {
	if deg > 360 || deg < 0 {
		return "wind direction outside of 0-360"
	}
	if deg >= 348.75 || deg < 11.25 {
		return "↑↑"
	} // N
	if deg >= 11.25 && deg < 33.75 {
		return "↑↗"
	} // NNE
	if deg >= 33.75 && deg < 56.25 {
		return "↗↗"
	} // NE
	if deg >= 56.25 && deg < 78.75 {
		return "→↗"
	} // ENE
	if deg >= 78.75 && deg < 101.25 {
		return "→→"
	} // E
	if deg >= 101.25 && deg < 123.75 {
		return "→↘"
	} // ESE
	if deg >= 123.75 && deg < 146.25 {
		return "↘↘"
	} // SE
	if deg >= 146.25 && deg < 168.75 {
		return "↓↘"
	} // SSE
	if deg >= 168.75 && deg < 191.25 {
		return "↓↓"
	} // S
	if deg >= 191.25 && deg < 213.75 {
		return "↓↙"
	} // SSW
	if deg >= 213.75 && deg < 236.25 {
		return "↙↙"
	} // SW
	if deg >= 236.25 && deg < 258.75 {
		return "←↙"
	} // WSW
	if deg >= 258.75 && deg < 281.25 {
		return "←←"
	} // W
	if deg >= 281.25 && deg < 303.75 {
		return "←↖"
	} // WNW
	if deg >= 303.75 && deg < 326.25 {
		return "↖↖"
	} // NW
	if deg >= 326.25 && deg < 348.75 {
		return "↑↖"
	} // NNW
	return "ERR: wind direction ranges do not overlap"
}

func colorizeTemperature(temperature float64) string {
	var temperatureColor string
	switch {
	case temperature < -10:
		temperatureColor = "#4444ff"
		break
	case temperature < 0:
		temperatureColor = "#9999ff"
		break
	case temperature < 5:
		temperatureColor = "#bbbbff"
		break
	case temperature < 10:
		temperatureColor = "#cccccc"
		break
	case temperature < 15:
		temperatureColor = "#00aa00"
		break
	case temperature < 25:
		temperatureColor = "#00dd00"
		break
	case temperature < 30:
		temperatureColor = "#bbaa00"
		break
	case temperature < 35:
		temperatureColor = "#aa4400"
		break
	case temperature < 40:
		temperatureColor = "#aa0000"
		break
	case temperature < 60:
		temperatureColor = "#ff0000"
		break
	default: // we're on sun, yaay
		temperatureColor = "#aa00aa"
		break
	}
	return fmt.Sprintf(`<span color="%s">%2.2f</span>`, temperatureColor, temperature)
}
