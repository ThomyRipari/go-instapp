package api

import (
	"log"
	"net/http"

	userHandlers "github.com/ThomyRipari/go-instapp/handlers"
	"github.com/ThomyRipari/go-instapp/middlewares"
	"github.com/ThomyRipari/go-instapp/services"
	"github.com/gorilla/mux"
)

func InitServer() {
	router := createRouter()

	server := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
	}

	log.Fatal(server.ListenAndServe())
}

func createRouter() *mux.Router {
	r := mux.NewRouter()

	client, err := services.ConnectMongoDB()

	if err != nil {
		log.Fatal(err)
	}
	log.Print("Succesfully conected to MongoDB Cloud Service")

	subRouter := r.PathPrefix("/api/v1").Subrouter()

	subRouter.HandleFunc("/register", middlewares.Logging(userHandlers.Register(client))).Methods("POST")

	return r
}
