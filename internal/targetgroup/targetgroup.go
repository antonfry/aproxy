package targetgroup

import (
	"aproxy/internal/backend"
	healchcheck "aproxy/internal/healthcheck"
	"errors"
	"fmt"
)

var ErrUnknownProto = errors.New("unknown proto")

type TargetGroup struct {
	Backends    []*backend.UDPBackend    `json:"servers"`
	Port        string                   `json:"port"`
	Proto       string                   `json:"proto"`
	HealthCheck *healchcheck.HealthCheck `json:"healthcheck"`
}

func New(config *Config, hcheck *healchcheck.HealthCheck) *TargetGroup {
	backends := make([]*backend.UDPBackend, 0)
	for _, b := range config.Backends {
		backends = append(backends, backend.New(b, config.Port))
	}
	return &TargetGroup{
		Backends:    backends,
		Port:        config.Port,
		Proto:       config.Proto,
		HealthCheck: hcheck,
	}
}

func (tg *TargetGroup) String() string {
	return fmt.Sprintf("Backends: %v, Port: %s, Proto: %s", tg.Backends, tg.Port, tg.Proto)
}

func (tg *TargetGroup) Check() error {
	var checkErrors error
	for _, s := range tg.Backends {
		var err error
		switch tg.HealthCheck.Proto {
		case "TCP":
			err = tg.HealthCheck.TCPcheck(s.Host)
		case "HTTP":
			err = tg.HealthCheck.HTTPcheck(s.Host)
		default:
			return ErrUnknownProto
		}
		if err != nil {
			checkErrors = errors.Join(checkErrors, err)
		}
		if s.IsAlive() && err != nil {
			s.SetDead()
		} else if !s.IsAlive() && err == nil {
			s.SetAlive()
		}
	}
	return checkErrors
}
