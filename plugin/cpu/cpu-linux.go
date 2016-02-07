package cpu

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// Get total cpu ticks (with no argument) or ticks for cpuid (counting from 0)

func GetCpuTicks(cpuid ...int) (cpuStats []cpuTicks, err error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return cpuStats, errors.New(fmt.Sprintf("cant open /proc/stats: %s", err))
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	fields := strings.Fields(scanner.Text())
	for strings.Contains(fields[0],`cpu`) {
		ticks,err := parseCpuTicks(fields)
		if err != nil {
			return cpuStats, err
		}
		scanner.Scan()
		fields = strings.Fields(scanner.Text())
		cpuStats = append(cpuStats,ticks)
	}
	return cpuStats,err

}
func parseCpuTicks(fields []string) (ticks cpuTicks, err error) {
	if len(fields) < 5 || !strings.Contains(fields[0], "cpu") {
		return ticks, errors.New(fmt.Sprintf("cpu id over total cpu count, cant decode %+v", fields))
	}
	// man proc
	ticks.user, _ = strconv.ParseUint(fields[1], 10, 64)
	ticks.nice, _ = strconv.ParseUint(fields[2], 10, 64)
	ticks.system, _ = strconv.ParseUint(fields[3], 10, 64)
	ticks.idle, _ = strconv.ParseUint(fields[4], 10, 64)
	if len(fields) > 5 {
		ticks.iowait, _ = strconv.ParseUint(fields[5], 10, 64)
	}
	if len(fields) > 6 {
		ticks.irq, _ = strconv.ParseUint(fields[6], 10, 64)
	}
	if len(fields) > 7 {
		ticks.softirq, _ = strconv.ParseUint(fields[7], 10, 64)
	}
	if len(fields) > 8 {
		ticks.steal, _ = strconv.ParseUint(fields[8], 10, 64)
	}
	if len(fields) > 9 {
		ticks.guest, _ = strconv.ParseUint(fields[9], 10, 64)
	}
	if len(fields) > 10 {
		ticks.guestNice, _ = strconv.ParseUint(fields[10], 10, 64)
	}
	ticks.total = ticks.user + ticks.nice + ticks.system + ticks.idle + ticks.iowait + ticks.irq + ticks.softirq + ticks.steal + ticks.guest + ticks.guestNice
	return ticks, err
}

func GetCpuCount() int {
	return runtime.NumCPU()
}
