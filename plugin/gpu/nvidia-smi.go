package gpu

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

var reClock = regexp.MustCompile(`([\d\.]+)\s*(\S+)`)
var rePower = regexp.MustCompile(`([\d\.]+)\s*W`)

type nvSmiQuery struct {
	nvSmiLog smiLog `xml:nvidia_smi_log`
}

type smiLog struct {
	DriverVersion string  `xml:"driver_version"`
	GPUs          []nvGPU `xml:"gpu"`
}
type nvGPU struct {
	ID            string          `xml:"id,attr"`
	ProductName   string          `xml:"product_name"`
	PowerReadings nvPowerReadings `xml:"power_readings"`
	MaxClocks     nvMaxClocks     `xml:"max_clocks"`
}

type nvPowerReadings struct {
	PowerLimit string `xml:"power_limit"`
}
type nvMaxClocks struct {
	Graphics string `xml:"graphics_clock"`
	// StreamingMultiprocessor
	SM     string `xml:"sm_clock"`
	Memory string `xml:"mem_clock"`
	Video  string `xml:"video_clock"`
}

func parseSmiQuery(q []byte) (*smiLog, error) {
	var s smiLog
	err := xml.Unmarshal([]byte(q), &s)
	return &s, err
}

func nvGetMaximumValues(gpuid string) (g gpuInfo, err error) {
	cmd := exec.Command(`nvidia-smi`, `-q`, `-i`, gpuid, `-x`)
	var out bytes.Buffer
	var errout bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errout
	err = cmd.Run()
	if err != nil {
		return g, fmt.Errorf("couldn't get data from nvidia-smi: %s | %s", err, errout.String())
	}
	smiData, err := parseSmiQuery(out.Bytes())
	if err != nil {
		return g, fmt.Errorf("error parsing nvidia-smi data %s", err)
	}
	if len(smiData.GPUs) < 1 {
		return g, fmt.Errorf("GPU with id %s not found", gpuid)
	}
	// for some reason every other max is readable via csv except video, read it from xml
	videoClock := reClock.FindStringSubmatch(smiData.GPUs[0].MaxClocks.Video)
	if len(videoClock) > 1 {
		i, err := strconv.Atoi(videoClock[1])
		if err != nil {
			return g, fmt.Errorf("could not convert %+v to int: %s", videoClock, err)
		}
		g.ClockVideo = i
	}

	return g, nil

}
