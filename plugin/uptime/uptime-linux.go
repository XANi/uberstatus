package uptime

import (
	"os"
	"bufio"
	"strings"
	"time"
	"sync"
)

type uptimeReader struct {
	uptimeFile *os.File
	uptimeFileScanner *bufio.Scanner
	sync.Mutex
}



func (p *plugin) getUptime() time.Duration {
	r := &p.uptimeReader
	if r.uptimeFile == nil {
		file, err := os.Open("/proc/uptime")
		if err != nil {
			log.Warningf("Can't open /proc/uptime: %s", err)
			return (time.Second* (60 * 66 + 6))
		} else {
			r.Lock()
			r.uptimeFile = file
			r.uptimeFileScanner = bufio.NewScanner(file)
			r.Unlock()
		}
	}
	_, err := r.uptimeFile.Seek(0, 0)
	if err != nil {
		log.Errorf("/proc/uptime seek failed: %s", err)
		return (time.Second* (60 * 55 + 5))
	}
	r.uptimeFileScanner.Scan()
	fields := strings.Fields(r.uptimeFileScanner.Text())
	duration, err := time.ParseDuration(fields[0]+"s")
	_  = err
	return duration

}
