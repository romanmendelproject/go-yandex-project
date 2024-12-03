package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var cfg Config

type Config struct {
	App struct {
		MaxCPU   int    `yaml:"max_cpu"`
		LogLevel string `yaml:"log_level"`
	} `yaml:"app"`
	DB struct {
		DBIP  string `yaml:"db_ip"`
		DNDSN string `yaml:"db_dsn"`
		Key   string `yaml:"key"`
	} `yaml:"db"`
}

func ReadConfig() Config {
	f, err := os.Open("/home/user/practicum/go-yandex-project/configuration.yml")
	if err != nil {
		log.Error(err)
		panic(err)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return cfg
}
