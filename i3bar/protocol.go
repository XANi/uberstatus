package i3bar

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	//
	"github.com/XANi/uberstatus/uber"
)

type Header struct {
	Version int8 `json:"version"`
	// not used yet
	// Stop_signal  uint8 `json:"stop_signal"`
	// Cont_signal  uint8 `json:"cont_signal"`
	ClickEvents bool `json:"click_events"`
}

// http://i3wm.org/docs/i3bar-protocol.html
type Msg struct {
	FullText            string `json:"full_text"`                       // full text, when shortening is not required
	ShortText           string `json:"short_text,omitempty"`            // shortened version of text to use when bar is full
	Color               string `json:"color,omitempty"`                 // color in #ffff00
	BorderColor         string `json:"border,omitempty"`                // color in #ffff00
	BackgroundColor     string `json:"background,omitempty"`            // color in #ffff00
	Markup              string `json:"markup,omitempty"`                // markup, pango or none (default
	MinWidth            string `json:"min_width,omitempty"`             // width in pixels, or string which will be measured for min_width
	Align               string `json:"align,omitempty"`                 // left/right/center align when size of text is smaller than minWidth
	Name                string `json:"name,omitempty"`                  // block name (ignored by i3bar, but will be returned in event)
	Instance            string `json:"instance,omitempty"`              // block instance (ignored by i3bar, but will be returned in event)
	Urgent              bool   `json:"urgent,omitempty"`                // urgent flag
	Separator           bool   `json:"separator,omitempty"`             // draw eparator
	SeparatorBlockWidth int16  `json:"separator_block_width,omitempty"` //number of pixe
}

// incoming event
type Event struct {
	Name     string `json:"name"`
	Instance string `json:"instance"`
	Button   int    `json:"button"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
}

func NewEvent() (r Event) {
	return r
}

func NewMsg() (r Msg) {
	// return msg with defaults
	r.FullText = "?"
	r.Color = `#aaaaaa`
	r.Separator = true
	r.Align = `center`
	r.SeparatorBlockWidth = 9
	return r
}

func CreateMsg(update uber.Update) (r Msg) {
	msg := NewMsg()
	msg.Name = update.Name
	msg.Instance = update.Instance
	msg.FullText = update.FullText
	msg.ShortText = update.ShortText
	msg.Color = update.Color
	msg.BackgroundColor = update.BackgroundColor
	msg.BorderColor = update.BorderColor
	msg.Markup = update.Markup
	msg.Urgent = update.Urgent
	return msg
}

func NewHeader() (r Header) {
	r.Version = 1
	r.ClickEvents = true
	return r
}

func (r Msg) Encode() []byte {
	s, _ := json.Marshal(r)
	return s
}

func (r Header) Encode() []byte {
	s, _ := json.Marshal(r)
	return s
}

// raw events are emitted like this:
//
//     [
//     {"name":"","button":1,"x":3775,"y":7}
//     ,{"name":"","button":1,"x":3775,"y":7}
//     ,{"name":"","button":1,"x":3775,"y":7}
//     ,{"name":"","button":1,"x":3775,"y":7}
//     ,{"name":"","button":1,"x":3775,"y":7}

// so we need some preprocessing to make it work with json parser

func FilterRawEvent(in []byte) []byte {
	if in[0] == byte('[') {
		return []byte("{}\n")
	}
	re := regexp.MustCompile(`\,{`)
	return re.ReplaceAllLiteral(in, []byte(`{`))
}

func EventReader() chan uber.Event {
	queue := make(chan uber.Event, 1)
	go eventReaderLoop(queue)
	return queue
}

func eventReaderLoop(events chan uber.Event) {
	stdin := bufio.NewReader(os.Stdin)
	var ct int64
	for {
		m := NewEvent()
		line, _ := stdin.ReadBytes('\n')
		if len(line) == 0 {
			continue
		}
		json.Unmarshal(FilterRawEvent(line), &m)
		if m.Button == 0 {
			continue
		} // Button is always >0, if it isnt present we got crap. TODO This should probably be logged if it shows up too often
		// This conversion is a lil bit of a waste but thanks to that there is no need to "taint" main uber.Event with any json tags
		out := uber.Event{
			Name:     m.Name,
			Instance: m.Instance,
			Button:   m.Button,
			X:        m.X,
			Y:        m.Y,
		}
		select {
		case events <- out:
		default:
			t := time.Now().Unix()
			if t != ct {
				fmt.Fprintf(os.Stderr, "input channel full, discarding input!")
				ct = t
			}
		}

	}
}
