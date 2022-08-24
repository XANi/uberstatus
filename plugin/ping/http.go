package ping

import (
	"net/http"
	"net/url"
	"time"
)

var pingHttpClient = http.Client{
	Transport:     nil,
	CheckRedirect: nil,
	Jar:           nil,
	Timeout:       time.Second * 12,
}

func httpPing(addr *url.URL) pingResult {
	var ping pingResult
	timeStart := time.Now()
	_, err := pingHttpClient.Head(addr.String())
	timeEnd := time.Now()
	if err == nil {
		ping.Duration = timeEnd.Sub(timeStart)
		ping.Ok = true
	} else {
		ping.Ok = false
	}
	return ping
}
