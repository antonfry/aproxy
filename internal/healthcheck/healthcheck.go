package healthcheck

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

var ErrHTTPNot200 = errors.New("http status code is not 200")

type HealthCheck struct {
	Proto    string        `json:"proto"`
	Port     string        `json:"port"`
	URI      string        `json:"uri"`
	Interval time.Duration `json:"interval"`
	Timeout  time.Duration `json:"timeout"`
}

func New(config *Config) *HealthCheck {
	return &HealthCheck{
		Proto:    config.Proto,
		Port:     config.Port,
		URI:      config.URI,
		Interval: time.Duration(config.Interval * int(time.Second)),
		Timeout:  time.Duration(config.Timeout * int(time.Second)),
	}
}

func (h *HealthCheck) String() string {
	return fmt.Sprintf("Proto: %s, Port: %s, Interval: %v, Timeout: %v", h.Proto, h.Port, h.Interval, h.Timeout)
}

func (h *HealthCheck) TCPcheck(host string) error {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, h.Port), h.Timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func (h *HealthCheck) HTTPcheck(host string) error {
	addr := net.JoinHostPort(host, h.Port)
	healthCheckUrl, err := url.JoinPath("http://", addr, h.URI)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodGet, healthCheckUrl, nil)
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout: h.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ErrHTTPNot200
	}
	return nil
}
