package logging

import (
	"github.com/romanmendelproject/go-yandex-project/internal/client/config"

	log "github.com/sirupsen/logrus"
)

func SetLogLevel() {
	cfg := config.ReadConfig()
	switch cfg.App.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.Warning("Log level incorrect. Set level info")
		log.SetLevel(log.InfoLevel)
	}
}
