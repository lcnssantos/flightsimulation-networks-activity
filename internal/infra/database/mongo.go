package database

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDriver struct {
	connectionString string
	dbName           string
	client           *mongo.Client
}

func NewMongoDriver(connectionString string, dbName string) MongoDriver {
	return MongoDriver{connectionString: connectionString, dbName: dbName}
}

func (d MongoDriver) GetClient() (*mongo.Client, error) {
	if d.client != nil {
		log.Info().Msg("MongoDB client already exists")
		return d.client, nil
	}

	log.Info().Msg("Creating MongoDB client")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(20)*time.Second)

	defer cancel()

	client, err := mongo.
		Connect(ctx, options.
			Client().
			ApplyURI(d.connectionString).
			SetMaxConnecting(3).
			SetMaxPoolSize(3).
			SetMaxConnIdleTime(10*time.Second).
			SetMaxConnecting(3),
		)

	if err != nil {
		return nil, err
	}

	d.client = client

	return client, nil
}

func (d MongoDriver) GetDatabase() (*mongo.Database, error) {
	client, err := d.GetClient()

	if err != nil {
		return nil, err
	}

	return client.Database(d.dbName), nil
}
