package ConnectMongo

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(collectionName string) (*mongo.Database, *mongo.Collection, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := godotenv.Load("env/local/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Retrieve the MongoDB connection URL from environment variables.
	mongoURL := os.Getenv("MONGO_CLIENT_URL")
	if mongoURL == "" {
		return nil, nil, fmt.Errorf("MONGO_CLIENT_URL not set in environment variables")
	}

	// Connect to the MongoDB client.
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Check if the database "CodeCatalyst" exists.
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list database names: %w", err)
	}

	dbExists := false
	for _, dbName := range databases {
		if dbName == "EventManagement" {
			dbExists = true
			break
		}
	}

	db := client.Database(os.Getenv("DB_NAME"))

	// Initialize collections.
	if !dbExists {
		// Create default collections if the database does not exist.
		defaultCollections := []string{
			"admin", "superAdmin", "userData", "events",
		}

		for _, coll := range defaultCollections {
			if err := db.CreateCollection(ctx, coll); err != nil {
				return nil, nil, fmt.Errorf("failed to create collection '%s': %w", coll, err)
			}
		}
	}

	// Check if the specified collection exists.
	collectionExists := false
	collectionNames, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list collection names: %w", err)
	}

	for _, name := range collectionNames {
		if name == collectionName {
			collectionExists = true
			break
		}
	}

	// Create the collection if it does not exist.
	if !collectionExists {
		if err := db.CreateCollection(ctx, collectionName); err != nil {
			return nil, nil, fmt.Errorf("failed to create collection '%s': %w", collectionName, err)
		}
	}

	col := db.Collection(collectionName)
	return db, col, nil
}

func FindOneDoc(collectionName string, filter bson.M, projection bson.M) (*mongo.SingleResult, error) {

	// Connect to the MongoDB collection.
	db, col, err := ConnectMongo(collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	defer func() {
		if disconnectErr := db.Client().Disconnect(context.TODO()); disconnectErr != nil {
			fmt.Println("Error disconnecting MongoDB client:", disconnectErr)
		}
	}()

	// Create a context with a timeout to avoid long-running operations.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set up the find options for the query.
	findOptions := options.FindOne().SetProjection(projection)

	// Execute the query.
	result := col.FindOne(ctx, filter, findOptions)
	return result, nil
}

func InsertOneDoc(collection string, data interface{}) (*mongo.InsertOneResult, error) {
	// Connect to MongoDB with a timeout context to prevent hanging connections.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Establish the connection.
	db, col, err := ConnectMongo(collection)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ensure the client is disconnected once the operation completes.
	defer func() {
		if err := db.Client().Disconnect(ctx); err != nil {
			fmt.Println("Error disconnecting MongoDB client:", err)
		}
	}()

	// Insert the document into the collection.
	result, err := col.InsertOne(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %w", err)
	}

	return result, nil
}
