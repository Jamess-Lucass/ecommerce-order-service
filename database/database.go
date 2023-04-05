package database

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var ctx = context.Background()

func Connect(log *zap.Logger) *mongo.Database {
	server := os.Getenv("DB_HOST")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Sugar().Fatalf("Could not parse PORT to an integar: %v", err)
	}
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	u := &url.URL{
		Scheme: "mongodb",
		User:   url.UserPassword(username, password),
		Host:   fmt.Sprintf("%s:%d", server, port),
	}

	clientOptions := options.Client().ApplyURI(u.String())
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Sugar().Fatalf("error connecting to database: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Sugar().Fatalf("error pinging redis database: %v", err)
	}

	return client.Database("order")
}
