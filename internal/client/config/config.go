// Модуль для объявления конфигурации клиента
package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	cfg Config
)

type Config struct {
	App struct {
		LogLevel string `yaml:"log_level"`
	} `yaml:"app"`
	DB struct {
		DBIP  string `yaml:"db_ip"`
		DNDSN string `yaml:"db_dsn"`
		Key   string `yaml:"key"`
	} `yaml:"db"`
}

func ReadConfig() Config {
	f, err := os.Open("./configuration.yml")
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
