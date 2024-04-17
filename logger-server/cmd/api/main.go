package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Client struct {
}

func main() {
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func () {
		if err = mongoClient.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}
}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)

	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting to MongoDB", err)
		return nil, err
	}

	return client, nil
}
