package sender

import (
	"net/http"
	"os"

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

func (sender *Sender) SetToken(req *http.Request) error {
	token, err := os.ReadFile("/tmp/data")
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:   "Token",
		Value:  string(token),
		MaxAge: 300,
	}
	req.AddCookie(cookie)

	return nil
}
