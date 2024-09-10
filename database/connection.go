package database

import (
	"context"
	"fmt"
	"os"
	"qrgen/service/utilities"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongoDB() (*mongo.Database, error) {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("DB_CONNECTION_STRING")).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	ctx := context.TODO()

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	databseConn := client.Database("students")

	collectionNames, err := databseConn.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to list collection names: %v", err)
	}

	if !utilities.Contains(collectionNames, "class10") {
		err = databseConn.CreateCollection(ctx, "class10")
		if err != nil {
			return nil, fmt.Errorf("failed to create 'class10' collection: %v", err)
		}
	}

	if !utilities.Contains(collectionNames, "admin") {
		err = databseConn.CreateCollection(ctx, "admin")
		if err != nil {
			return nil, fmt.Errorf("failed to create 'admin' collection: %v", err)
		}
	}

	log.Info().Msg("pinged your deployment. You successfully connected to MongoDB!")

	return databseConn, nil
}
