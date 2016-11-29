package ping

import (
	"time"
	"net/http"
)

func httpPing(addr string, out chan *pingResult) {
	var okCount uint64
	var failCount uint64
	for {
		var ping pingResult
		timeStart := time.Now()
		_, err := http.Head(addr)
		timeEnd := time.Now()
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
