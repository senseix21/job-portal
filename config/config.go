package config

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func Connect() {
	// Set client options
	// mongoURI := "mongodb+srv://job-portal:;Rq^,y=2mzG=U+^@cluster0.t32indv.mongodb.net/job-portal?retryWrites=true&w=majority&appName=Cluster0"

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb+srv://job-portal:;Rq^,y=2mzG=U+^@cluster0.t32indv.mongodb.net/job-portal?retryWrites=true&w=majority&appName=Cluster0")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Ping the database to confirm connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	DB = client
	fmt.Println("Connected to MongoDB successfully!")
}

func GetCollection(database, collection string) *mongo.Collection {
	if DB == nil {
		log.Fatal("Database connection is not initialized")
	}
	return DB.Database(database).Collection(collection)
}
