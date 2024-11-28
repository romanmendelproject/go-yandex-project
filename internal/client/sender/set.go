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

func (sender *Sender) SetCred(name, username, password, meta string) error {
	var requestData types.CredType
	var requestBody bytes.Buffer

	requestData.Name = name
	requestData.Username = username
	requestData.Password = password
	requestData.Meta = meta

	body, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/api/user/update/cred/", sender.cfg.DB.DBIP)

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
			return nil
		}
	}

	return err
}
