package sender

import (
	"net/http"
	"os"
	"strconv"

	"github.com/romanmendelproject/go-yandex-project/internal/client/config"
	"github.com/romanmendelproject/go-yandex-project/internal/client/flags"
	log "github.com/sirupsen/logrus"
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

func (s *Sender) SetData() error {
	switch flags.DataType {
	case "cred":
		err := s.SetCred(flags.Name, flags.Username, flags.Password, flags.Meta)
		if err != nil {
			return err
		}
	case "text":
		err := s.SetText(flags.Name, flags.Data, flags.Meta)
		if err != nil {
			return err
		}
	case "byte":
		data := []byte(flags.Data)
		err := s.SetByte(flags.Name, data, flags.Meta)
		if err != nil {
			return err
		}
	case "card":
		data, err := strconv.ParseInt(flags.Data, 10, 64)
		if err != nil {
			log.Errorf("Data is not int value")
		}
		err = s.SetCard(flags.Name, data, flags.Meta)
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
	case "text":
		err := sender.GetText(flags.Name)
		if err != nil {
			return err
		}
	case "byte":
		err := sender.GetByte(flags.Name)
		if err != nil {
			return err
		}
	case "card":
		err := sender.GetCard(flags.Name)
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
