package userHandlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	models "github.com/ThomyRipari/go-instapp/types"
	"github.com/thedevsaddam/govalidator"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var user models.User

		err := decoder.Decode(&user)

		if err != nil {
			log.Print(err)

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rules := govalidator.MapData{
			"username": []string{"required", "between:3,15"},
			"email":    []string{"required", "min:4", "max:32", "email"},
		}

		opts := govalidator.Options{
			Data:  &user,
			Rules: rules,
		}

		instanceValidator := govalidator.New(opts)

		validationErr := instanceValidator.ValidateStruct()

		if len(validationErr) > 0 {
			data, _ := json.MarshalIndent(validationErr, "", "  ")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(data)
			return
		}

		usersCollection := client.Database("instapp_db").Collection("users")

		ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)

		insertedDoc, err := usersCollection.InsertOne(ctx, user)

		if err != nil {
			log.Print(err)

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		objectIDJSON, _ := insertedDoc.InsertedID.(primitive.ObjectID).MarshalJSON()

		w.WriteHeader(http.StatusCreated)
		w.Write(objectIDJSON)
	}
}
