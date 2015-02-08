package main

import (
	"encoding/json"
	"fmt"
	"i3bar"
	"os"
	"plugin"
	"plugin_interface"
	"regexp"
	"time"
)

func main() {
	button := 0
	header := i3bar.NewHeader()
	msg := i3bar.NewMsg()
	msg.FullText = `test`
	b, err := json.Marshal(header)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	fmt.Println("\n[")
	//c, err := json.Marshal(msg)

	c := msg.Encode()

	i3input := i3bar.EventReader()

	updates := make(chan plugin_interface.Update, 10)

	_ = plugin.NewPlugin("clock", "", c, updates)

	for {
		fmt.Print(`[`)
		msg := i3bar.NewMsg()
		msg.FullText = fmt.Sprintf("Btn: %d", button)
		os.Stdout.Write(msg.Encode())
		fmt.Print(`,`)
		os.Stdout.Write(c)
		fmt.Print(`,`)
		os.Stdout.Write(getTime())
		fmt.Println(`],`)

		select {
		case ev := (<-i3input):
			button = ev.Button
		case upd := <-updates:
			c = i3bar.CreateMsg(upd).Encode()
		case <-time.After(time.Second):
			button = 0
		}
	}
}

func getTime() []byte {
	msg := i3bar.NewMsg()
	msg.Name = "clock"
	t := time.Now().Local()
	// reference Mon Jan 2 15:04:05 MST 2006 (unix: 1136239445)
	msg.FullText = t.Format(`15:04:05`)
	msg.Color = `#ffffff`
	return msg.Encode()
}

func San(in []byte) []byte {
	re := regexp.MustCompile(`\,{`)
	return re.ReplaceAllLiteral(in, []byte(`{`))
}
