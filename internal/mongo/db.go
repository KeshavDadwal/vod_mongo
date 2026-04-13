package mongo

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Config struct {
	URI    string
	DBName string
}

func ConfigFromEnv() (Config, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return Config{}, errors.New("missing MONGODB_URI")
	}
	dbName := os.Getenv("MONGODB_DB")
	if dbName == "" {
		dbName = "tiktok"
	}
	return Config{URI: uri, DBName: dbName}, nil
}

func Connect(ctx context.Context, cfg Config) (*mongo.Client, *mongo.Database, error) {
	clientOpts := options.Client().
		ApplyURI(cfg.URI).
		SetRetryWrites(true)

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, nil, err
	}

	return client, client.Database(cfg.DBName), nil
}

