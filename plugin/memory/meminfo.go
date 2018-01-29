package memory

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
)

// Format:
// MemTotal:       16362772 kB
// MemFree:          222548 kB
// MemAvailable:   12783960 kB
// Buffers:            5668 kB
// Cached:         11981816 kB
// SwapCached:            0 kB
// ...
// SwapTotal:       7815584 kB
// SwapFree:        7815584 kB
// ...
// struct has all in (or converted to) bytes
type memInfo struct {
	Total      int64
	Free       int64
	Buffers    int64
	Cached     int64
	Available  int64 // since 3.14
	HasAvailable bool
	SwapTotal  int64
	SwapFree   int64
	SwapCached int64
}

var memRegex = regexp.MustCompile(`^(\S+)\:\s+(\d+)\skB`)

func getMemInfo() (mem memInfo) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		match := memRegex.FindStringSubmatch(line)
		if len(match) < 3 {
			continue
		}
		i, err := strconv.ParseInt(match[2], 10, 64)
		if err != nil {
			log.Warning(`cant convert %s from match %s to int`, match[2], match[1])
		}
		i = i * 1024
		switch {
		case match[1] == `MemTotal`:
			mem.Total = i
		case match[1] == `MemFree`:
			mem.Free = i
		case match[1] == `Buffers`:
			mem.Buffers = i
		case match[1] == `Cached`:
			mem.Cached = i
		case match[1] == `MemAvailable`:
			mem.Available = i
			mem.HasAvailable = true
		case match[1] == `SwapCached`:
			mem.SwapCached = i
		case match[1] == `SwapFree`:
			mem.SwapFree = i
		case match[1] == `SwapTotal`:
			mem.SwapTotal = i
		}
	}
	return mem
}
