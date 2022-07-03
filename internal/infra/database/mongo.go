package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDriver struct {
	connectionString string
	dbName           string
}

func NewMongoDriver(connectionString string, dbName string) MongoDriver {
	return MongoDriver{connectionString: connectionString, dbName: dbName}
}

func (d MongoDriver) GetClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(20)*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.connectionString))

	if err != nil {
		return nil, err
	}

	return client, nil
}

func (d MongoDriver) GetDatabase() (*mongo.Database, error) {
	client, err := d.GetClient()
	if err != nil {
		return nil, err
	}

	return client.Database(d.dbName), nil
}
