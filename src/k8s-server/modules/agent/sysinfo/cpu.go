package sysinfo

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type CPU struct {
	Core    int
	Threads int
	Cache   int64
	Speed   int64
	CPU     int
	Model   string
	Vendor  string
}

const (
	CPUInfoFile = "/proc/cpuinfo"
)

var (
	reCPUSeparator = regexp.MustCompile("\t+: ")
	reCPUSpace     = regexp.MustCompile(" +")
	reCPUCacheSize = regexp.MustCompile(`^(\d+) KB$`)
)

func (c *Collector) getCPUInfo() {
	f, err := os.Open(CPUInfoFile)
	if err != nil {
		return
	}
	defer f.Close()

	cpu := make(map[string]bool)
	core := make(map[string]bool)

	var cpuID string

	b := bufio.NewScanner(f)
	for b.Scan() {
		if col := reCPUSeparator.Split(b.Text(), 2); col != nil {
			switch col[0] {
			case "vendor_id":
				if c.Info.CPU.Vendor == "" {
					c.Info.CPU.Vendor = col[1]
				}
			case "model name":
				if c.Info.CPU.Model == "" {
					model := reCPUSpace.ReplaceAllLiteralString(col[1], " ")
					c.Info.CPU.Model = strings.Replace(model, "- ", "-", 1)
				}
			case "cpu MHz":
				if c.Info.CPU.Speed == 0 {
					if speed, err := strconv.ParseFloat(col[1], 64); err == nil {
						speed = math.Ceil(speed)
						c.Info.CPU.Speed = int64(speed)
					}
				}
			case "physical id":
				cpuID = col[1]
				cpu[cpuID] = true
			case "core id":
				coreID := fmt.Sprintf("%s%s", cpuID, col[1])
				core[coreID] = true
			case "cache size":
				if c.Info.CPU.Cache == 0 {
					if cache := reCPUCacheSize.FindStringSubmatch(col[1]); cache != nil {
						if _cache, err := strconv.ParseInt(cache[1], 10, 64); err == nil {
							c.Info.CPU.Cache = _cache
						}

					}
				}
			}
		}
	}
	if b.Err() != nil {
		return
	}

	c.Info.CPU.CPU = len(cpu)
	c.Info.CPU.Core = len(core)

	c.Info.CPU.Threads = runtime.NumCPU()
}
