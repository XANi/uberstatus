package ping

import (
	"time"
	"net"
)

func tcpPing(addr string, out chan *pingResult,t time.Duration) {
	var okCount uint64
	var failCount uint64
	for {
		var ping pingResult
		timeStart := time.Now()
		dial := net.Dialer{Timeout: time.Duration(time.Second * 10)}
		c, err := dial.Dial("tcp", addr)
		timeEnd := time.Now()
		if err == nil {
			c.Close()
			okCount = okCount + 1
			ping.Duration = timeEnd.Sub(timeStart)
			ping.Ok = true
		} else {
			failCount = failCount + 1
			ping.Ok = false
		}
		ping.OkCount = okCount
		ping.FailCount = failCount
		out <- &ping
		time.Sleep(t)
	}

}
