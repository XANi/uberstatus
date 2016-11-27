package ping

import (
	"time"
	"net"
)

func tcpPing(addr string, out chan *pingResult) {
	var okCount uint64
	var failCount uint64
	for {
		var ping pingResult
		timeStart := time.Now()
		c, err := net.Dial("tcp", addr)
		timeEnd := time.Now()
		c.Close()
		if err == nil {
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
	}

}
