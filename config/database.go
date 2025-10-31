package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB establishes a connection to MongoDB
// It connects to MongoDB running on localhost:27017
// and verifies the connection with a ping
// Returns: *mongo.Client - The MongoDB client instance
func ConnectDB() *mongo.Client {
	log.Println("Attempting to connect to MongoDB...")
	// Configure MongoDB connection URI
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Create context with timeout for connection attempt
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Attempt to connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Verify the connection is alive with a ping
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB is not reachable: %v", err)
	}

	log.Println("Successfully connected to MongoDB!")
	return client
}

// Client is the global MongoDB client instance
// Initialized when the package is loaded
var Client *mongo.Client = ConnectDB()

// OpenCollection returns a reference to a MongoDB collection
// Parameters:
//   - collectionName: Name of the collection to access
//
// Returns:
//   - *mongo.Collection: Reference to the requested collection in the 'usersdb' database
func OpenCollection(collectionName string) *mongo.Collection {

	if Client == nil {
		log.Fatal("MongoDB client is not initialized. Please call ConnectDB first.")
	}
	// Return the specified collection from the 'usersdb' database
	return Client.Database("usersdb").Collection(collectionName)
}
