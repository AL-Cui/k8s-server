package sysinfo

import (
	"fmt"
	"k8s-server/utils/logs"
	"net"
)

type Network struct {
	Name       string
	IP         string
	MACAddress string
}

func (c *Collector) getNetworkInfo() error {
	c.Info.Networks = []Network{}
	interfaces, err := net.Interfaces()
	if err != nil {
		logs.Error("get network info failed: %v", err)
		return err
	}
	for _, i := range interfaces {
		if i.Name == "lo" {
			continue
		}
		addrs, err := i.Addrs()
		if err != nil {
			logs.Error("get network addr failed: %v", err)
			return err
		}
		if len(addrs) == 0 {
			continue
		}
		network := Network{
			Name:       i.Name,
			IP:         addrs[0].String(),
			MACAddress: i.HardwareAddr.String(),
		}
		isSame, ipAddr, err := c.isSameNetworkSegment(network.IP)
		if err != nil {
			logs.Error("can't report whether the network includes ip: %v", err)
			return err
		}
		// collect nodes IP and MAC info
		if isSame && ipAddr != "" && network.MACAddress != "" {
			conn := c.redisPool.Get()
			defer conn.Close()
			key := fmt.Sprintf("nodes:%s", c.hostname)
			conn.Do("HMSET", key, "IP", ipAddr, "MAC", network.MACAddress)
		}
		c.Info.Networks = append(c.Info.Networks, network)
	}
	return nil
}

// if the incoming IP is in the same network segment with mgmt IP, it returns
// true.
func (c *Collector) isSameNetworkSegment(ipStr string) (bool, string, error) {
	ip, _, err := net.ParseCIDR(ipStr)
	if err != nil {
		return false, "", err
	}
	return c.mgmtIPNet.Contains(ip), ip.String(), nil
}
