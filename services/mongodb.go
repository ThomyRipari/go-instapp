package services

import (
	"context"
	"log"
	"time"

	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"labix.org/v2/mgo/bson"
)

func ConnectMongoDB() ([]string, error) {
	dbUser := ""
	dbPassword := ""
	dbHost := "@cluster0-u9by7.mongodb.net/test?retryWrites=true&w=majority"
	mongoURI := fmt.Sprintf("mongodb+srv://%s:%s%s", dbUser, dbPassword, dbHost)

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	if err != nil {
		log.Print(err)

		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	connErr := client.Connect(ctx)

	if connErr != nil {
		log.Print(err)

		return nil, connErr
	}

	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Print(err)

		return nil, err
	}

	databases, err := client.ListDatabaseNames(ctx, bson.M{})

	if err != nil {
		log.Print(err)

		return nil, err
	}

	log.Fatal(databases)

	return databases, nil
}
