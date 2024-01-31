package conf

import (
	"aproxy/internal/aproxy"
	"aproxy/internal/healthcheck"
	"aproxy/internal/targetgroup"
)

type Config struct {
	Server      aproxy.Config      `yaml:"server"`
	Targetgroup targetgroup.Config `yaml:"targetgroup"`
	Healthcheck healthcheck.Config `yaml:"healthcheck"`
}
