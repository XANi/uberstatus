package gpu

import (
    "encoding/xml"
)


type nvSmiQuery struct {
    nvSmiLog smiLog `xml:nvidia_smi_log`
}

type smiLog struct {
    DriverVersion string `xml:"driver_version"`
    GPUs []nvGPU `xml:"gpu"`
}
type nvGPU struct {
    ID string `xml:"id,attr"`
    ProductName string `xml:"product_name"`
    PowerReadings nvPowerReadings `xml:"power_readings"`
    MaxClocks nvMaxClocks `xml:"max_clocks"`
}

type nvPowerReadings struct {
    PowerLimit string `xml:"power_limit"`
}
type nvMaxClocks struct {
    Graphics string `xml:"graphics_clock"`
    // StreamingMultiprocessor
    SM string `xml:"sm_clock"`
    Memory string `xml:"mem_clock"`
    Video string `xml:"video_clock"`
}

func parseSmiQuery(q []byte) (*smiLog, error) {
    var s smiLog
    err := xml.Unmarshal([]byte(q), &s)
    return &s, err

}