package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongo connects to MongoDB and returns client and selected database.
func ConnectMongo(ctx context.Context, uri string, dbName string) (*mongo.Client, *mongo.Database, error) {
	clientOpts := options.Client().ApplyURI(uri)

	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		return nil, nil, err
	}

	// connect with timeout context (caller supplies ctx)
	if err := client.Connect(ctx); err != nil {
		return nil, nil, err
	}

	// ping
	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, nil, err
	}

	db := client.Database(dbName)
	return client, db, nil
}
