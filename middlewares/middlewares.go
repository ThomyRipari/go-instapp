package middlewares

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	models "github.com/ThomyRipari/go-instapp/types"
	"github.com/dgrijalva/jwt-go"
)

func Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date := time.Now()

		defer func() {
			log.Println(r.URL.Path, time.Since(date))
		}()

		next(w, r)
	}
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		stringToken := strings.Split(authHeader, " ")

		var claims models.Claims

		token, err := jwt.ParseWithClaims(stringToken[1], &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("A0Zr98j/3yX R~XHH!jmN]LWX/,?RT"), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				log.Print(err)

				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !token.Valid {
			log.Print("Token no valido")
			w.WriteHeader(http.StatusUnauthorized)

			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, models.ContextTokenClaims, claims)

		next(w, r.WithContext(ctx))
	}
}
