package config

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

// ConnectDB establishes a connection to MongoDB.
func ConnectDB() (*mongo.Client, error) {
	// Use an environment variable for sensitive data
	mongoURI := "mongodb+srv://SoeelBeg:SOEELbeg@cluster0.qmoic.mongodb.net/?retryWrites=true&w=majority"
	clientOptions := options.Client().ApplyURI(mongoURI) // MongoDB URI
	client, err := mongo.NewClient(clientOptions)        // Create MongoDB client
	if err != nil {
		return nil, fmt.Errorf("failed to create Mongo client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx) // Connect to MongoDB
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil) // Ping MongoDB
	if err != nil {
		return nil, fmt.Errorf("MongoDB ping failed: %v", err)
	}

	fmt.Println("Connected to MongoDB successfully!")
	DB = client
	return client, nil
}

// GetCollection returns a MongoDB collection.
func GetCollection(dbName, collectionName string) *mongo.Collection {
	return DB.Database(dbName).Collection(collectionName)
}
