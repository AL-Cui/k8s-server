package logs

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"os"

	"github.com/gomodule/redigo/redis"
)

// redisWriter implements LoggerInterface.
// it writes messages to a redis channel .
type redisWriter struct {
	Key               string `json:"key"`
	Addr              string `json:"addr"`
	HostName          string `json:"host"`
	Port              int    `json:"port"`
	Password          string `json:"password"`
	DataType          string `json:"data_type"`
	Timeout           int    `json:"timeout"`
	ReconnectInterval int    `json:"reconnect_interval"`
	Level             int    `json:"level"`
	conns             *redis.Pool
}

// newRedisLogger creates new redisWriter returning as LoggerInterface.
func newRedisLogger() Logger {
	return &redisWriter{
		Level:    LevelTrace,
		Key:      "hpc-logs",
		DataType: "channel",
		Port:     6379,
		Addr:     "127.0.0.1",
	}
}

// Init initializes redisWriter with json config.
func (r *redisWriter) Init(jsonConfig string) error {
	err := json.Unmarshal([]byte(jsonConfig), r)
	if err != nil {
		return err
	}
	if r.HostName == "" {
		hostName, err := os.Hostname()
		if err != nil {
			return fmt.Errorf("logger can't get correct hostname:%v", err)
		}
		r.HostName = hostName
	}
	r.conns = initRedisPool(fmt.Sprintf("%s:%d", r.Addr, r.Port), r.Password)
	return err
}

// WriteMsg write message in connection.
// if connection is down, try to re-connect.
func (r *redisWriter) WriteMsg(when time.Time, msg string, level int) (err error) {
	conn := r.conns.Get()
	defer conn.Close()
	h, _ := formatTimeHeader(when)
	msg = fmt.Sprintf("[%s] %s %s", r.HostName, string(h), msg)
	switch r.DataType {
	case "list":
		_, err = conn.Do("RPUSH", r.Key, msg)
	case "set":
		_, err = conn.Do("SET", r.Key, msg)
	case "channel":
		_, err = conn.Do("PUBLISH", r.Key, msg)
	default:
		err = errors.New("unknown DataType: " + r.DataType)
	}
	return
}

// init a redis connection pool.
func initRedisPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     200 * 2,
		MaxActive:   200 * 2,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password == "" {
				return c, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

// Destroy is empty.
func (r *redisWriter) Destroy() {
}

// Flush is empty.
func (r *redisWriter) Flush() {
}

func init() {
	Register(AdapterRedis, newRedisLogger)
}
