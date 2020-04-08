package middlewares

import (
	"time"
	"net/http"
	"log"
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