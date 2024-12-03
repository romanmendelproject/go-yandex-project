package sender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/romanmendelproject/go-yandex-project/internal/crypto"
	"github.com/romanmendelproject/go-yandex-project/internal/types"

	log "github.com/sirupsen/logrus"
)

func (s *Sender) SetCred(name, username, password, meta string) error {
	var requestData types.CredType

	requestData.Name = name
	requestData.Username = username
	requestData.Password = password
	requestData.Meta = meta

	body, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/api/user/update/cred/", s.cfg.DB.DBIP)

	err = s.sendData(url, body)
	if err != nil {
		return err
	}

	return err
}

func (s *Sender) SetText(name, data, meta string) error {
	var requestData types.TextType

	requestData.Name = name
	requestData.Data = data
	requestData.Meta = meta

	body, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/api/user/update/text/", s.cfg.DB.DBIP)

	s.sendData(url, body)
	if err != nil {
		return err
	}

	return err
}

func (s *Sender) SetByte(name string, data []byte, meta string) error {
	var requestData types.ByteType

	requestData.Name = name
	requestData.Data = data
	requestData.Meta = meta

	body, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/api/user/update/byte/", s.cfg.DB.DBIP)

	s.sendData(url, body)
	if err != nil {
		return err
	}

	return err
}

func (s *Sender) SetCard(name string, data int64, meta string) error {
	var requestData types.CardType

	requestData.Name = name
	requestData.Data = data
	requestData.Meta = meta

	body, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/api/user/update/card/", s.cfg.DB.DBIP)

	s.sendData(url, body)
	if err != nil {
		return err
	}

	return err
}

func (sender *Sender) sendData(url string, body []byte) error {
	var requestBody bytes.Buffer

	gz := gzip.NewWriter(&requestBody)
	gz.Write(body)
	gz.Close()

	client := http.Client{}

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		log.Error(err)
	}

	err = sender.SetToken(req)
	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	if sender.cfg.DB.Key != "" {
		hash := crypto.GetHash(body, sender.cfg.DB.Key)
		req.Header.Set("HashSHA256", hash)
	}

	for _, timeSleep := range retries {
		resp, err := client.Do(req)
		if err != nil {
			log.Errorf("Failed to send to server: %s. Retrying after %ds...", err, timeSleep)
			time.Sleep(time.Duration(timeSleep) * time.Second)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return fmt.Errorf("not expected status code: %d", resp.StatusCode)
		} else {
			println("Data saved successfully")
			return nil
		}
	}

	return nil
}
