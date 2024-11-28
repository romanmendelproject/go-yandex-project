package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/romanmendelproject/go-yandex-project/internal/types"
	log "github.com/sirupsen/logrus"
)

func (sender *Sender) getRequest(name string) (*http.Response, error) {
	url := fmt.Sprintf("http://%s/api/user/value/cred/%s", sender.cfg.DB.DBIP, name)
	var requestBody bytes.Buffer
	client := http.Client{}

	req, err := http.NewRequest("GET", url, &requestBody)
	if err != nil {
		log.Error(err)
	}

	err = sender.SetToken(req)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (sender *Sender) GetCred(name string) error {
	var request types.CredType
	resp, err := sender.getRequest(name)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(&request); err != nil {
		return err
	}

	fmt.Printf("Username: %s", request.Username)
	fmt.Printf("Password: %s", request.Password)
	fmt.Printf("Meta: %s", request.Meta)

	return nil
}
