package ping

import (
	"time"
	"net"
)

func tcpPing(addr string) pingResult {
	var ping pingResult
	timeStart := time.Now()
	dial := net.Dialer{Timeout: time.Duration(time.Second * 10)}
	c, err := dial.Dial("tcp", addr)
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
