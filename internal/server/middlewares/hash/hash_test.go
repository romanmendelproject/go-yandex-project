package hash

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/romanmendelproject/go-yandex-project/internal/crypto"
	"github.com/stretchr/testify/assert"
)

func TestHashMiddleware(t *testing.T) {
	key := "secret"

	t.Run("successful hash verification", func(t *testing.T) {

		metrics := []byte("test metrics")
		expectedHash := crypto.GetHash(metrics, key)

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := HashMiddleware(key)(testHandler)

		req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(metrics))
		req.Header.Set("HashSHA256", expectedHash)

		recorder := httptest.NewRecorder()

		middleware.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("missing hash header", func(t *testing.T) {
		metrics := []byte("test metrics")

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := HashMiddleware(key)(testHandler)

		req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(metrics))

		recorder := httptest.NewRecorder()

		middleware.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("invalid hash", func(t *testing.T) {
		metrics := []byte("test metrics")

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := HashMiddleware(key)(testHandler)

		req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(metrics))

		req.Header.Set("HashSHA256", "invalidhash")

		recorder := httptest.NewRecorder()

		middleware.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}
