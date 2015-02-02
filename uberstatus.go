package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"i3bar"
	"os"
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

	st := bufio.NewReader(os.Stdin)
	f, _ := os.Create("/tmp/dat2")
	w := bufio.NewWriter(f)
	//	dec := json.NewDecoder(st)
	for {
		time.Sleep(time.Second)
		fmt.Print(`[`)
		msg := i3bar.NewMsg()
		msg.FullText = fmt.Sprintf("Btn: %d", button)
		os.Stdout.Write(msg.Encode())
		fmt.Print(`,`)
		os.Stdout.Write(c)
		fmt.Print(`,`)
		os.Stdout.Write(getTime())
		fmt.Println(`],`)
		//		m := i3bar.NewEvent()
		//		dec.Decode(&m)
		str, _ := st.ReadBytes('\n')
		w.Write(i3bar.FilterRawEvent(str))
		w.Flush()
		//		button = m.Button
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
