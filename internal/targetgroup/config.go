package targetgroup

type Config struct {
	Backends []string `yaml:"backends"`
	Port     string   `yaml:"port"`
	Proto    string   `yaml:"proto"`
}

func NewConfig() *Config {
	return &Config{
		Backends: []string{"172.16.51.192", "172.16.51.220", "172.16.51.189", "172.16.51.191", "172.16.51.190", "172.16.51.201"},
		Port:     "9600",
		Proto:    "UDP",
	}
}
