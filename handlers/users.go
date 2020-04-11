package userHandlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	models "github.com/ThomyRipari/go-instapp/types"
	"github.com/dgrijalva/jwt-go"
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
			"firstname": []string{"required", "between:3,32"},
			"surname":   []string{"required", "between:3,32"},
			"username":  []string{"required", "between:3,20"},
			"email":     []string{"required", "min:4", "max:32", "email"},
			"password":  []string{"required", "between:8,32"},
			"age":       []string{"numeric_between:18,105"},
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(objectIDJSON)
	}
}

func Login(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials models.Credentials
		json.NewDecoder(r.Body).Decode(&credentials)

		usersCollection := client.Database("instapp_db").Collection("users")

		ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)

		var docResult bson.M
		err := usersCollection.FindOne(ctx, bson.D{{"username", credentials.Username}}).Decode(&docResult)

		if err != nil {
			log.Print(err)

			if err == mongo.ErrNoDocuments {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("The username or password is incorrect"))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if credentials.Password != docResult["password"] {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("The username or password is incorrect"))
			return
		}

		expirationTime := time.Now().Add(1 * time.Minute)

		claims := models.Claims{
			Username: credentials.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte("A0Zr98j/3yX R~XHH!jmN]LWX/,?RT"))

		if err != nil {
			log.Print(err)

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userDocToJSON, err := json.Marshal(docResult)

		if err != nil {
			log.Print(err)

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Authorization", "Bearer "+tokenString)
		w.Write(userDocToJSON)
		w.WriteHeader(http.StatusOK)
	}
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	interfaceClaims := r.Context().Value(models.ContextTokenClaims)
	claims := interfaceClaims.(models.Claims)

	if time.Unix(claims.StandardClaims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		log.Print("El token todavia tiene validez")

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expirationTime := time.Now().Add(1 * time.Minute)

	claims.StandardClaims.ExpiresAt = expirationTime.Unix()

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := refreshToken.SignedString([]byte("A0Zr98j/3yX R~XHH!jmN]LWX/,?RT"))

	if err != nil {
		log.Print(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+tokenString)
	w.WriteHeader(http.StatusOK)
}

func Social(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("Social Handler")
	}
}
