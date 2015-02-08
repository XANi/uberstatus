package i3bar

import (
	"encoding/json"
	"regexp"
	"bufio"
	"os"
	"plugin_interface"
//	"fmt"
)

type I3barHeader struct{
	Version int8 `json:"version"`
	// not used yet
	// Stop_signal  uint8 `json:"stop_signal"`
	// Cont_signal  uint8 `json:"cont_signal"`
	ClickEvents bool `json:"click_events"`
}
// http://i3wm.org/docs/i3bar-protocol.html
type I3barMsg struct {
	FullText string `json:"full_text"` // full text, when shortening is not required
	ShortText string `json:"short_text"` // shortened version of text to use when bar is full
	Color string `json:"color"` // color in #ffff00
	MinWidth string `json:"min_width"` // width in pixels, or string which will be measured for min_width
	Align string `json:"align"` // left/right/center align when size of text is smaller than minWidth
	Name string `json:"name"` // block name (ignored by i3bar, but will be returned in event)
	Instance string `json:"instace"` // block instance (ignored by i3bar, but will be returned in event)
	Urgent bool `json:"urgent"` // urgent flag
	Separator bool `json:"separator"` // draw eparator
	SeparatorBlockWidth int16 `json:"separator_block_width"` //number of pixe
}

// incoming event
type I3barEvent struct {
	Name string `json:"name"`
	Instance string `json:"instance"`
	Button int `json:"button"`
	X int `json:"x"`
	Y int `json:"y"`
}

func NewEvent() (r I3barEvent) {
	return r
}

func NewMsg() (r I3barMsg) {
	// return msg with defaults
	r.FullText="asd"
	r.Color=`#aaaaaa`
	r.Separator = true
	r.Align = `center`
	r.SeparatorBlockWidth = 9
	return r
}

func CreateMsg(update plugin_interface.Update) (r I3barMsg) {
	msg := NewMsg()
	msg.Name = update.Name
	msg.Instance = update.Instance
	msg.FullText = update.FullText
	msg.ShortText = update.ShortText
	msg.Color = update.Color
	msg.Urgent = update.Urgent
	return msg
}

func NewHeader() (r I3barHeader) {
	r.Version = 1
	r.ClickEvents = true
	return r
}

func (r I3barMsg) Encode() []byte {
	s, _ :=  json.Marshal(r)
	return s
}

func (r I3barHeader) Encode() []byte {
	s, _ :=  json.Marshal(r)
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

func EventReader() chan I3barEvent {
	queue := make(chan I3barEvent)
	go eventReaderLoop(queue)
	return queue
}

func eventReaderLoop(events chan I3barEvent) {
	stdin := bufio.NewReader(os.Stdin)
	for {
		m := NewEvent()
		line, _ := stdin.ReadBytes('\n')
		json.Unmarshal( FilterRawEvent(line), &m)
		if(m.X == 0) { continue } // x is neccesary, if it isnt present we got crap. THis should probably be logged... once I figure out how to make app-wide logger...
		events <- m
	}
}
