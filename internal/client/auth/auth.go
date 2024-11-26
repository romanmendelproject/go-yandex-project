package auth

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/romanmendelproject/go-yandex-project/internal/client/config"
	"github.com/romanmendelproject/go-yandex-project/internal/crypto"

	log "github.com/sirupsen/logrus"
)

var retries = []int{1, 3, 5}

type AuthData struct {
	Login    string
	Password string
}

type UserAuth struct {
	cfg config.Config
}

func NewUserAuth(cfg config.Config) *UserAuth {
	return &UserAuth{
		cfg: cfg,
	}
}

func (auth *UserAuth) Register(login, password string) error {
	var requestData AuthData
	var requestBody bytes.Buffer

	requestData.Login = login
	requestData.Password = password

	body, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/api/user/register", auth.cfg.DB.DBIP)

	gz := gzip.NewWriter(&requestBody)
	gz.Write(body)
	gz.Close()

	client := http.Client{}

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		log.Error(err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	if auth.cfg.DB.Key != "" {
		hash := crypto.GetHash(body, auth.cfg.DB.Key)
		req.Header.Set("HashSHA256", hash)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("not expected status code: %d", resp.StatusCode)
	}

	return nil

}

func (auth *UserAuth) Login(login, password string) error {
	var requestData AuthData
	var requestBody bytes.Buffer

	requestData.Login = login
	requestData.Password = password

	body, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/api/user/login", auth.cfg.DB.DBIP)

	gz := gzip.NewWriter(&requestBody)
	gz.Write(body)
	gz.Close()

	client := http.Client{}

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		log.Error(err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	if auth.cfg.DB.Key != "" {
		hash := crypto.GetHash(body, auth.cfg.DB.Key)
		req.Header.Set("HashSHA256", hash)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	token := getToken(resp)
	if token == "" {
		return fmt.Errorf("token is empty")
	}

	f, err := os.Create("/tmp/data")
	defer f.Close()
	_, err = f.WriteString(token)

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("not expected status code: %d", resp.StatusCode)
	}

	println("Auth Success")

	return nil

}

func getToken(resp *http.Response) string {
	cookies := resp.Cookies()
	for _, c := range cookies {
		if c.Name == "Token" {
			return c.Value
		}
	}
	return ""
}
