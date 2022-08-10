package ping

import (
	"net/http"
	"time"
)

var pingHttpClient = http.Client{
	Transport:     nil,
	CheckRedirect: nil,
	Jar:           nil,
	Timeout:       time.Second * 12,
}

func httpPing(addr string) pingResult {
	var ping pingResult
	timeStart := time.Now()
	_, err := pingHttpClient.Head(addr)
	timeEnd := time.Now()
	if err == nil {
		ping.Duration = timeEnd.Sub(timeStart)
		ping.Ok = true
	} else {
		ping.Ok = false
	}
	return ping
}
