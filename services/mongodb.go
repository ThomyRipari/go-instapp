package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectMongoDB() (*mongo.Client, error) {
	log.Print("Connecting to Mongo DB Cloud Service...")

	clientOpts := options.Client().ApplyURI("mongodb+srv://:@cluster0-u9by7.mongodb.net")

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	client, err := mongo.Connect(ctx, clientOpts)

	if err != nil {
		log.Print(err)

		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Print(err)

		return nil, err
	}

	return client, nil
}
