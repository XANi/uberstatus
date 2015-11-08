package uber


// event passed from X/status bar manager
type Event struct{
	Name string
	Instance string
	Button int
	X int
	Y int
}

// update sent by plugin
type Update struct{
	Name string
	Instance string
	FullText string  // full text, when shortening is not required
	ShortText string  // shortened version of text to use when bar is full
	Color string // color in #ffff00
	Urgent bool `json:"urgent"` // urgent flag, will update (assuming backend allows) immediately if that flag is present
}

type Tag struct {
	Name string
	Instance string
}

type TaggedUpdate struct {
	Update *Update
	Tag *Tag
}
