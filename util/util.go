package util

import (
	"fmt"
	"text/template"
)

// calculate divider and unit for bytes
func GetUnitBytes(bytes int64) (divider int, unit string) {
	switch {
	case bytes < 10000:
		return 1, ``
	case bytes < 2000*1024:
		return 1024, `K`
	case bytes < 5000*1024*1024:
		return 1024 * 1024, `M`
	default:
		return 1024 * 1024 * 1024, `G`
	}
}

// format bytes
func FormatUnitBytes(bytes int64) (s string) {
	div, unit := GetUnitBytes(bytes)
	return fmt.Sprintf("%4.2f%+1s", float64(bytes)/float64(div), unit)
}

// generate bar chart from percent
func GetBarChar(pct int) string {
	switch {
	case pct > 90:
		return `█`
	case pct > 80:
		return `▇`
	case pct > 70:
		return `▆`
	case pct > 60:
		return `▅`
	case pct > 40:
		return `▄`
	case pct > 20:
		return `▂`
	case pct > 10:
		return `▁`
	}
	return ` `
}

// generate color from percentage (0 - good/green 100 - bad/red)
func GetColorPct(pct int) string {
	switch {
	case pct > 90:
		return `#dd0000`
	case pct > 80:
		return `#cc3333`
	case pct > 70:
		return `#ccaa44`
	case pct > 50:
		return `#cc9966`
	case pct > 30:
		return `#cccc66`
	case pct > 15:
		return `#66cc66`
	}
	return `#666666`
}

func NewTemplate(name string, tpl string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"percentToColor": GetColorPct,
		"percentToBar": GetBarChar,
		"formatBytes": FormatUnitBytes,
	}
	return template.New(name).Funcs(funcMap).Parse(tpl)
}
