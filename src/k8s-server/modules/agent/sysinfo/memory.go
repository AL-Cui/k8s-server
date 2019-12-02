package sysinfo

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strconv"
)

type Memory struct {
	MemTotal  int
	MemFree   int
	SwapTotal int
	SwapFree  int
}

const (
	memInfoFile = "/proc/meminfo"
)

var (
	reMemSeparator = regexp.MustCompile(":[\t ]+")
	reMemSpace     = regexp.MustCompile(" +")
	reMemSize      = regexp.MustCompile(`^(\d+)`)
)

func parseSize(s string) (i int, err error) {
	if sz := reMemSize.FindStringSubmatch(s); sz != nil {
		if size, err := strconv.ParseInt(sz[1], 10, 64); err == nil {
			return int(size), nil
		}
		return i, err
	}
	return i, errors.New("no match string")
}

func (c *Collector) getMemoryInfo() {
	f, err := os.Open(memInfoFile)
	if err != nil {
		return
	}
	defer f.Close()

	b := bufio.NewScanner(f)
	for b.Scan() {
		if col := reMemSeparator.Split(b.Text(), 2); col != nil {
			switch col[0] {
			case "MemTotal":
				if c.Info.Memory.MemTotal == 0 {
					if memTotal, err := parseSize(col[1]); err == nil {
						c.Info.Memory.MemTotal = memTotal
					}
				}
			case "MemFree":
				if c.Info.Memory.MemFree == 0 {
					if memFree, err := parseSize(col[1]); err == nil {
						c.Info.Memory.MemFree = memFree
					}
				}
			case "SwapTotal":
				if c.Info.Memory.SwapTotal == 0 {
					if swapTotal, err := parseSize(col[1]); err == nil {
						c.Info.Memory.SwapTotal = swapTotal
					}
				}
			case "SwapFree":
				if c.Info.Memory.SwapFree == 0 {
					if swapFree, err := parseSize(col[1]); err == nil {
						c.Info.Memory.SwapFree = swapFree
					}
				}
			}
		}
	}
}
