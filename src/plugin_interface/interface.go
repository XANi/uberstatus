package plugin_interface

type Event struct{
	Button int
	X int
	Y int
}

type Update struct{
	FullText string  // full text, when shortening is not required
	ShortText string  // shortened version of text to use when bar is full
	Color string // color in #ffff00
	Urgent bool `json:"urgent"` // urgent flag, will update (assuming backend allows) immediately if that flag is present
}
