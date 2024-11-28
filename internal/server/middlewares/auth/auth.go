package auth

import (
	"net/http"
	"strconv"

	"github.com/romanmendelproject/go-yandex-project/internal/server/jwt"
	log "github.com/sirupsen/logrus"
)

func IsAuthorized(token *jwt.JWT) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqToken, err := r.Cookie("Token")
			if err != nil {
				log.Error("error getting token", "error", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)

				return
			}

			userID, err := token.ParseToken(reqToken.Value)
			if err != nil || userID == 0 {
				http.Error(w, err.Error(), http.StatusUnauthorized)

				return
			}

			r.Header.Set("userID", strconv.Itoa(userID))

			next.ServeHTTP(w, r)
		})
	}

}
