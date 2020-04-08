package api

import (
	"github.com/gorilla/mux"
	"github.com/ThomyRipari/go-instapp/handlers"
	"github.com/ThomyRipari/go-instapp/middlewares"
	"net/http"
	"log"
)

func InitServer() {
	router := createRouter()

	server := &http.Server{
		Addr: "127.0.0.1:8000",
	}

	http.Handle("/", router)

	log.Fatal(server.ListenAndServe())
}

func createRouter() *mux.Router {
	r := mux.NewRouter()

	sub_router := r.PathPrefix("/api/v1").Subrouter()

	sub_router.HandleFunc("/register", middlewares.Logging(userHandlers.Register)).Methods("POST")

	return r
}