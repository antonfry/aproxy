package targetgroup

type Config struct {
	Backends []string `yaml:"backends"`
	Port     string   `yaml:"port"`
	Proto    string   `yaml:"proto"`
}

func NewConfig() *Config {
	return &Config{
		Backends: []string{},
		Port:     "",
		Proto:    "",
	}
}
