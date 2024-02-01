package main

import (
	"aproxy/internal/aproxy"
	"aproxy/internal/conf"
	"aproxy/internal/healthcheck"
	"aproxy/internal/roundrobin"
	"aproxy/internal/targetgroup"
	"flag"
	"os"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
	configPath   string
)

func main() {
	flag.StringVar(&configPath, "conf", "/etc/aproxy.yml", "path to yml configuration")
	flag.Parse()
	config := mustConf(configPath)

	healthCheck := healthcheck.New(&config.Healthcheck)
	targetGroup := targetgroup.New(&config.Targetgroup, healthCheck)
	pool := roundrobin.New(targetGroup)
	server := aproxy.New(&config.Server, pool)

	log.Info("Build version: ", buildVersion)
	log.Info("Build date: ", buildDate)
	log.Info("Build commit: ", buildCommit)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func mustConf(configPath string) conf.Config {
	var config conf.Config
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		log.WithError(err).Fatal("Unable to read configuration: ", configPath)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.WithError(err).Fatal("Unable to Unmarshal configuration: ", configPath)
	}
	return config
}
