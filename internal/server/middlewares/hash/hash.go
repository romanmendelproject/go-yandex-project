package hash

import (
	"bytes"
	"io"
	"net/http"

	"github.com/romanmendelproject/go-yandex-project/internal/crypto"

	log "github.com/sirupsen/logrus"
)

// HashMiddleware декодирует запросы к серверу
func HashMiddleware(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		logFn := func(res http.ResponseWriter, req *http.Request) {
			metrics, err := io.ReadAll(req.Body)
			if err != nil {
				log.Error(err)
			}

			hash := req.Header.Get("HashSHA256")
			if hash == "" {
				log.Error("Missing hash header")
				res.WriteHeader(http.StatusBadRequest)

				return
			}

			expectedHash := crypto.GetHash(metrics, key)

			if hash != expectedHash {
				log.Error("Hash is not valid")

				res.WriteHeader(http.StatusBadRequest)
			}

			req.Body = io.NopCloser(bytes.NewBuffer(metrics))

			next.ServeHTTP(res, req)
		}
		return http.HandlerFunc(logFn)
	}
}
