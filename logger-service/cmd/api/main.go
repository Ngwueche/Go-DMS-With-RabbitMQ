package main

import (
	"context"
	"log"
	"logger/cmd/api/data"
	"net/http"
	"os"
	"time"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const webPort = "8080"

type Application struct {
	DB     *mongo.Client
	Models data.Models
}

func main() {
	log.Println("Starting logger service...")
	db := connectToMongoDB()
	if db == nil {
		log.Fatal("Failed to connect to mongo database")
	}
	app := Application{
		DB:     db,
		Models: data.New(db),
	}

	server := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func connectToMongoDB() *mongo.Client {
	uri := os.Getenv("MONGO_URI")

	if uri == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}
	for range 10 {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		client, err := mongo.Connect(options.Client().ApplyURI(uri))
		if err == nil {
			err = client.Ping(ctx, nil)
			cancel()

			if err == nil {
				log.Println("Connected to MongoDB!")
				return client
			}
		}
		cancel()
		log.Println("Waiting for Mongo")
		time.Sleep(2 * time.Second)
	}
	return nil

}
