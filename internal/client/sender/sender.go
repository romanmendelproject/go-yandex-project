package sender

import (
	"github.com/romanmendelproject/go-yandex-project/internal/client/config"
	"github.com/romanmendelproject/go-yandex-project/internal/client/flags"
)

var retries = []int{1, 3, 5}

type Sender struct {
	cfg config.Config
}

func NewSender(cfg config.Config) *Sender {
	return &Sender{
		cfg: cfg,
	}
}

func (sender *Sender) SetData() error {
	switch flags.DataType {
	case "cred":
		err := sender.SetCred(flags.Name, flags.Username, flags.Password, flags.Meta)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sender *Sender) GetData() error {
	switch flags.DataType {
	case "cred":
		err := sender.GetCred(flags.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
