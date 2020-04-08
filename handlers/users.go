package userHandlers

import (
	"net/http"
	// "fmt"
	"github.com/ThomyRipari/go-instapp/services"
)

func Register(w http.ResponseWriter, r *http.Request) {
	_, err := services.ConnectMongoDB()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}