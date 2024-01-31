package healthcheck

type Config struct {
	Proto    string `yaml:"proto"`
	Port     string `yaml:"port"`
	Interval int    `yaml:"interval"`
	Timeout  int    `yaml:"timeout"`
	URI      string `yaml:"uri"`
}

func NewConfig() *Config {
	return &Config{
		Proto:    "HTTP",
		Port:     "9600",
		Interval: 10,
		Timeout:  1,
		URI:      "",
	}
}
