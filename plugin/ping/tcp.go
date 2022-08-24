package ping

import (
	"net"
	"net/url"
	"time"
)

func tcpPing(addr *url.URL) pingResult {
	var ping pingResult
	timeStart := time.Now()
	dial := net.Dialer{Timeout: time.Duration(time.Second * 10)}
	c, err := dial.Dial("tcp", addr.Host)
	timeEnd := time.Now()
	if err == nil {
		c.Close()
		ping.Duration = timeEnd.Sub(timeStart)
		ping.Ok = true
	} else {
		ping.Ok = false
	}
	return ping
}
