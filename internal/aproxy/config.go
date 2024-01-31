package aproxy

import "fmt"

type Config struct {
	Workers  int    `yaml:"workers"`
	LogLevel string `yaml:"loglevel"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Buffer   int    `yaml:"buffer"`
}

func NewConfig() *Config {
	return &Config{
		LogLevel: "INFO",
		Workers:  1,
		Host:     "0.0.0.0",
		Port:     "9600",
		Buffer:   2048,
	}
}

func (c *Config) String() string {
	return fmt.Sprintf("logLevel: %s, workers: %d, host: %s, port: %s, buffer: %d", c.LogLevel, c.Workers, c.Host, c.Port, c.Buffer)
}
