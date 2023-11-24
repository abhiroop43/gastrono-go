package database

import (
	"context"
	"fmt"
	"gastrono-go/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"time"
)

func DBInstance() *mongo.Client {

	f, err := os.Open("config.yml")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	var cfg models.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	MongoDb := cfg.Database.Connection
	fmt.Print(MongoDb)

	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Connected to MongoDB")
	return client
}

var Client *mongo.Client = DBInstance()


func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = Client.Database("gastrono-go").Collection(collectionName)
	return collection
}
