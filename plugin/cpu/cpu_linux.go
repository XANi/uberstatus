package cpu
import (
	"os"
	"bufio"
	"strings"
)


func GetCpuTicks() (ticks cpuTicks, err  error) {
	file, err := os.Open("/proc/stat")
    if err != nil {
        log.Fatal(err)
  		return ticks, err
    }
    defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	fields :=  strings.Split(scanner.Text()," ")
	_ = fields
	return ticks, err
}
