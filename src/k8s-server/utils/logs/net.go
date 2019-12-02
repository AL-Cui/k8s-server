package logs

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/fatih/pool"
)

// netWriter implements loggerInterface.
// it writes message to a tcp/udp channel.
type netWriter struct {
	Protocol string `json:"protocol"`
	Addr     string `json:"addr"`
	HostName string `json:"hostname"`
	Port     int    `json:"port"`
	Level    int    `json:"level"`
	Prefix   string `json:"prefix"`
	conns    pool.Pool
}

// newNetLogger creates new netWriter returning as LoggerInterface.
func newNetLogger() Logger {
	return &netWriter{
		Protocol: "udp",
		Addr:     "127.0.0.1",
		Port:     8888,
		Level:    LevelTrace,
		Prefix:   "LOG:",
	}
}

// Init initializes netWriter with json config.
func (n *netWriter) Init(jsonConfig string) error {
	err := json.Unmarshal([]byte(jsonConfig), n)
	if err != nil {
		return err
	}
	if n.HostName == "" {
		hostName, err := os.Hostname()
		if err != nil {
			return fmt.Errorf("logger can't get correct hostname:%v", err)
		}
		n.HostName = hostName
	}
	n.conns, err = initNetPool(n.Protocol, n.Addr, n.Port)
	if err != nil {
		return err
	}
	return nil
}

// WriteMsg write message in connection.
// if connection is down, try to re-connect.
func (n *netWriter) WriteMsg(when time.Time, msg string, level int) (err error) {
	conn, err := n.conns.Get()
	if err != nil {
		return err
	}
	defer conn.Close()
	h, _ := formatTimeHeader(when)
	msg = fmt.Sprintf("%s %s %s", n.Prefix, string(h), msg)
	_, err = conn.Write([]byte(msg))
	return err
}

// init a net connection pool.
func initNetPool(protocol, addr string, port int) (pool.Pool, error) {
	factory := func() (net.Conn, error) {
		return net.Dial(protocol, fmt.Sprintf("%s:%d", addr, port))
	}
	return pool.NewChannelPool(5, 30, factory)
}

// Destroy is empty.
func (n *netWriter) Destroy() {}

// Flush is empty.
func (n *netWriter) Flush() {}

func init() {
	Register(AdapterNet, newNetLogger)
}
