package clock

import (
	"plugin_interface"
	"time"
)

func New(config interface{}, events chan plugin_interface.Event, update chan plugin_interface.Update) {
	for {
		select {
		case _ = (<-events):
			Update(update)
		case <-time.After(time.Second):
			Update(update)
		}
	}

}

func Update(update chan plugin_interface.Update) {
	time := GetTimeEvent(`2006-01-02 MST 15:04:05`)
	time.Color=`#DDDDFF`
	update <- time
}


func GetTimeEvent(format string) plugin_interface.Update {
	t :=  time.Now().Local()
	var ev plugin_interface.Update
	ev.FullText = t.Format(format)
	ev.ShortText = t.Format(`15:04:05`)
	return ev
}
