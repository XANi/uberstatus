package cpu
import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"errors"
	"fmt"
)


// Get total cpu ticks (with no argument) or ticks for cpuid (counting from 0)

func GetCpuTicks(cpuid ...int) (ticks cpuTicks, err  error) {
	file, err := os.Open("/proc/stat")
    if err != nil {
  		return ticks, errors.New(fmt.Sprintf("cant open /proc/stats: %s", err))
    }
    defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	log.Debug("--- %+v", cpuid)
	if len(cpuid) > 0 {
		log.Debug("--- %+v", cpuid)
		for i := -1; i < cpuid[0]; i++ { scanner.Scan() }
	}
	fields :=  strings.Split(scanner.Text()," ")
	_ = fields
	if len(fields) < 5 || !strings.Contains(fields[0],"cpu") {
		return ticks, errors.New(fmt.Sprintf("cpu id over total cpu count, cant decode %+v", fields))
	}
	// man proc
	ticks.user, _ = strconv.ParseUint(fields[2],10,64)
	ticks.nice, _ = strconv.ParseUint(fields[3],10,64)
	ticks.system, _ = strconv.ParseUint(fields[4],10,64)
	ticks.idle, _ = strconv.ParseUint(fields[5],10,64)
	if len(fields) >= 6 { ticks.iowait, _ = strconv.ParseUint(fields[6],10,64) }
	if len(fields) >= 7 { ticks.irq, _ = strconv.ParseUint(fields[7],10,64) }
	if len(fields) >= 8 { ticks.softirq, _ = strconv.ParseUint(fields[8],10,64) }
	if len(fields) >= 9 { ticks.steal, _ = strconv.ParseUint(fields[9],10,64) }
	if len(fields) >= 10 { ticks.guest, _ = strconv.ParseUint(fields[10],10,64) }
	if len(fields) >= 11 { ticks.guestNice, _ = strconv.ParseUint(fields[11],10,64) }
	ticks.total = ticks.user + ticks.nice + ticks.system + ticks.idle + ticks.iowait + ticks.irq + ticks.softirq + ticks.steal + ticks.guest  + ticks.guestNice
	log.Debug("%+v", ticks)
	return ticks, err
}
