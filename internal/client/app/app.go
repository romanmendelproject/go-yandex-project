// Модуль для инициализации и запуска агента
package app

import (
	"github.com/romanmendelproject/go-yandex-project/internal/client/auth"
	"github.com/romanmendelproject/go-yandex-project/internal/client/config"
	"github.com/romanmendelproject/go-yandex-project/internal/client/flags"
	"github.com/romanmendelproject/go-yandex-project/internal/client/sender"

	log "github.com/sirupsen/logrus"
)

// ExecCommand запускает клиентскую программу
func ExecCommand() {
	flags.ParseFlags()
	cfg := config.ReadConfig()
	sender := sender.NewSender(cfg)

	switch flags.Action {
	case "set":
		err := sender.SetData()
		checkError(err)
	case "get":
		err := sender.GetData()
		checkError(err)
	case "reg":
		userAuth := auth.NewUserAuth(cfg)
		err := userAuth.Register(flags.Username, flags.Password)
		checkError(err)
	case "auth":
		userAuth := auth.NewUserAuth(cfg)
		err := userAuth.Login(flags.Username, flags.Password)
		checkError(err)
	default:
		panic("Error action key -a")
	}
}

func checkError(e error) {
	if e != nil {
		log.Error(e)
	}
}
