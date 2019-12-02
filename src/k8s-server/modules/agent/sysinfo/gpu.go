package sysinfo

import (
	"bytes"
	"encoding/xml"
	"os/exec"
)

const (
	cmdNvidiaSMI = "nvidia-smi"
)

//GPU contains complete GPU information
type GPU struct {
	DriverVersion string   `xml:"driver_version"`
	GPUs          []GPUXML `xml:"gpu"`
}

//GPUXML contains single GPU information from the XML file
type GPUXML struct {
	ProductName string          `xml:"product_name"`
	FanSpeed    string          `xml:"fan_speed"`
	FBMemory    FBMemoryInfo    `xml:"fb_memory_usage"`
	Temperature TemperatureInfo `xml:"temperature"`
	Clocks      ClocksInfo      `xml:"clocks"`
}

// FBMemoryInfo contains on-board frame buffer memory information.
type FBMemoryInfo struct {
	Total string `xml:"total"`
	Used  string `xml:"used"`
	Free  string `xml:"free"`
}

//TemperatureInfo contains readings from temperature sensors on the board.
type TemperatureInfo struct {
	GPUTemperature    string `xml:"gpu_temp"`
	MemoryTemperature string `xml:"memory_temp"`
}

//ClocksInfo contains current frequency at which parts of the GPU are running.
type ClocksInfo struct {
	GraphicsClock string `xml:"graphics_clock"`
	SMClock       string `xml:"sm_clock"`
	MEMClock      string `xml:"mem_clock"`
	VideoClock    string `xml:"video_clock"`
}

func (c *Collector) getGPUInfo() {
	output, err := exec.Command(cmdNvidiaSMI, "-q", "-x").Output()
	if err != nil {
		return
	}
	var data GPU
	r := bytes.NewReader(output)
	x := xml.NewDecoder(r)
	err = x.Decode(&data)
	if err != nil {
		return
	}
	c.Info.GPU = data
}
