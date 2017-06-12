package ping

import (
	"time"
	"net/http"
)

func httpPing(addr string) pingResult {
	var ping pingResult
	timeStart := time.Now()
	_, err := http.Head(addr)
	timeEnd := time.Now()
	if err == nil {
		ping.Duration = timeEnd.Sub(timeStart)
		ping.Ok = true
	} else {
		ping.Ok = false
	}
	return ping
}
