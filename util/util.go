package util

import (
	"fmt"
	"text/template"
	"bytes"
	"time"
	"html"
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
	case pct > 5:
		return `#669966`
	}
	return `#667766`
}

type Template struct{
	*template.Template
	buf *bytes.Buffer
}


func NewTemplate(name string, tpl string) (*Template, error) {
	funcMap := template.FuncMap{
		"percentToColor": GetColorPct,
		"percentToBar": GetBarChar,
		"formatBytes": FormatUnitBytes,
		"formatDuration": FormatDuration,
		"formatDurationPadded": FormatDurationPadded,
		"escape": html.EscapeString,
		"color": func (color string, text string) string{
			return `<span color="` + color + `">` + text + `</span>`
		},

	}
	t, err := template.New(name).Funcs(funcMap).Parse(tpl)
	return &Template{t, new(bytes.Buffer)}, err
}

func (t Template)ExecuteString(i interface{}) (string) {
	t.buf.Reset()
	err := t.Execute(t.buf,i)
	if err != nil {
		return fmt.Sprintf("tpl [%+v] error: %s",i, err)
	} else {
		return t.buf.String()
	}
}
// Formats duration. Does not go over 7 characters till 9999+h
func FormatDuration(t time.Duration) string {
	if t.Hours() > 1 {
		hours := int(t.Hours())
		minutes := int( int(t.Minutes()) - (hours * 60) )
		return fmt.Sprintf("%dh%0.2dm",hours,minutes)
	}
	if t.Hours() == 1 {
		return fmt.Sprintf("%5.0fm",t.Minutes())
	}
	if t.Minutes() > 1 {
		return fmt.Sprintf("%5.1fm",t.Minutes())
	}
	if t.Seconds() > 4 {
		return fmt.Sprintf("%5.2fs",t.Seconds())
	}
	if t.Seconds() >= 1 {
		return fmt.Sprintf("%5.0fms",t.Seconds() * 1000)
	}
	if t.Seconds() >= 0.1 {
		return fmt.Sprintf("%5.1fms",t.Seconds() * 1000)
	}
	if t.Seconds() > 0.001 {
		return fmt.Sprintf("%5.2fms",t.Seconds() * 1000)
	}
	if t.Nanoseconds() >= 100000 { //100us
		return fmt.Sprintf("%5.1fµs",float64(t.Nanoseconds())/1000)
	}
	if t.Nanoseconds() >= 10000 { // 10us
		return fmt.Sprintf("%5.2fµs",float64(t.Nanoseconds())/1000)
	}
	if t.Nanoseconds() >= 1000 { // 1us
		return fmt.Sprintf("%5.2fµs",float64(t.Nanoseconds())/1000)
	} else {
		return fmt.Sprintf("%5dns",t.Nanoseconds())
	}
}

// FormatDuration with added padding to 7 (
func FormatDurationPadded(t time.Duration) string {
	return fmt.Sprintf("%7s", FormatDuration(t))
}

func WaitForTs(nextTs *time.Time) {
	t := time.Now()
	for nextTs.After(t) {
		diff :=nextTs.Sub(t)
		// cap sleeping at 10s in case date changes between ticks
		if diff > time.Second * 10  {
			time.Sleep(time.Second * 10)
		} else {
			time.Sleep(diff)
		}
		t = time.Now()
	}
}

func GenerateColorBarLookupTable() (colorTable map[int8]string, barTable map[int8]string) {
	var i int8
	var ltBar = make(map[int8]string)
	var ltColor = make(map[int8]string)
	for i = -1; i <= 101; i++ {
		color := GetColorPct(int(i))
		ltColor[i] = color
		ltBar[i] = `<span color="` + color + `">` + GetBarChar(int(i)) + `</span>`
	}
	return ltColor, ltBar
}
