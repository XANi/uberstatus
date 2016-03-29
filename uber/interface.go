package uber

// event passed from X/status bar manager
type Event struct {
	Name     string
	Instance string
	Button   int
	X        int
	Y        int
}

// update sent by plugin
type Update struct {
	Name            string
	Instance        string
	FullText        string // full text, when shortening is not required
	ShortText       string // shortened version of text to use when bar is full
	Color           string // color in #ffff00
	BackgroundColor string // color in #ffff00
	BorderColor     string // color in #ffff00
	Markup          string // markup, so far only pango in i3bar is supported
	Urgent          bool   `json:"urgent"` // urgent flag, will update (assuming backend allows) immediately if that flag is present
}

type Tag struct {
	Name     string
	Instance string
}

type TaggedUpdate struct {
	Update *Update
	Tag    *Tag
}

type PluginConfig struct {
	Name     string
	Instance string
	Config   map[string]interface{}
	Events   chan Event
	Update   chan Update
}
