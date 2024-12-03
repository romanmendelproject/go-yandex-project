package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/romanmendelproject/go-yandex-project/internal/types"
	log "github.com/sirupsen/logrus"
)

func (s *Sender) getRequest(typeData, name string) (*http.Response, error) {
	url := fmt.Sprintf("http://%s/api/user/value/%s/%s", s.cfg.DB.DBIP, typeData, name)
	var requestBody bytes.Buffer
	client := http.Client{}

	req, err := http.NewRequest("GET", url, &requestBody)
	if err != nil {
		log.Error(err)
	}

	err = s.SetToken(req)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Sender) GetCred(name string) error {
	var request types.CredType
	resp, err := s.getRequest("cred", name)
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		return fmt.Errorf("data with the specified name was not found")
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error receiving data %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&request); err != nil {
		return err
	}

	fmt.Printf("Username: %s", request.Username)
	println()
	fmt.Printf("Password: %s", request.Password)
	println()
	fmt.Printf("Meta: %s", request.Meta)

	return nil
}

func (s *Sender) GetText(name string) error {
	var request types.TextType
	resp, err := s.getRequest("text", name)
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		return fmt.Errorf("data with the specified name was not found")
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error receiving data %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&request); err != nil {
		return err
	}

	data := strings.TrimSpace(request.Data)
	fmt.Printf("Data: %s", data)
	println()
	fmt.Printf("Meta: %s", request.Meta)

	return nil
}

func (s *Sender) GetByte(name string) error {
	var request types.ByteType
	resp, err := s.getRequest("byte", name)
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		return fmt.Errorf("data with the specified name was not found")
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error receiving data %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&request); err != nil {
		return err
	}

	fmt.Printf("Data: %s", request.Data)
	println()
	fmt.Printf("Meta: %s", request.Meta)

	return nil
}

func (s *Sender) GetCard(name string) error {
	var request types.CardType
	resp, err := s.getRequest("card", name)
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		return fmt.Errorf("data with the specified name was not found")
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error receiving data %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&request); err != nil {
		return err
	}

	fmt.Printf("Data: %d", request.Data)
	println()
	fmt.Printf("Meta: %s", request.Meta)

	return nil
}
