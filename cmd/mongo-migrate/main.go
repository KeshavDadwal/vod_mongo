package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	appmongo "github.com/vod/internal/mongo"
	"github.com/vod/routes"
)

const defaultPort = "8080"

func main() {
	// Load .env when present (safe to ignore if missing in production).
	_ = godotenv.Load()

	cfg, err := appmongo.ConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, db, err := appmongo.Connect(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	if err := appmongo.EnsureValidators(ctx, db); err != nil {
		log.Fatal(err)
	}
	if err := appmongo.EnsureIndexes(ctx, db); err != nil {
		log.Fatal(err)
	}

	log.Printf("mongo migration applied (db=%s)", cfg.DBName)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := routes.NewGraphQLRouter(db)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(router.Listen(":" + port))
}
