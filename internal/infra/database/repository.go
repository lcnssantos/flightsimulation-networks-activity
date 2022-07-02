package database

import (
	"context"
	"time"

	"github.com/lcnssantos/online-activity/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	driver MongoDriver
}

func NewRepository(driver MongoDriver) Repository {
	return Repository{driver: driver}
}

func (r *Repository) getMongoDB() (*mongo.Database, error) {
	return r.driver.GetDatabase()
}

func (r *Repository) SaveBrazilActivity(ctx context.Context, activity domain.NetworkActivity) error {
	db, err := r.getMongoDB()

	if err != nil {
		return err
	}

	if _, err := db.Collection("br_activity").InsertOne(ctx, &activity); err != nil {
		return err
	}

	return nil
}

func (r *Repository) SaveActivity(ctx context.Context, activity domain.NetworkActivity) error {
	db, err := r.getMongoDB()

	if err != nil {
		return err
	}

	if _, err := db.Collection("activity").InsertOne(ctx, &activity); err != nil {
		return err
	}

	return nil
}

func (r *Repository) SaveGeoActivity(ctx context.Context, activity domain.GeoNetworkActivity) error {
	db, err := r.getMongoDB()

	if err != nil {
		return err
	}

	if _, err := db.Collection("geo_activity").InsertOne(ctx, &activity); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetBrazilActivityByMinutes(ctx context.Context, minutes int64) ([]domain.NetworkActivity, error) {
	output := make([]domain.NetworkActivity, 0)

	db, err := r.getMongoDB()

	if err != nil {
		return output, err
	}

	cursor, err := db.Collection("br_activity").Find(ctx, bson.M{"date": bson.M{
		"$gte": primitive.NewDateTimeFromTime(time.Now().Add(-time.Duration(minutes) * time.Minute)),
	}})

	if err != nil {
		return output, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &output)

	return output, err
}

func (r *Repository) GetActivityByMinutes(ctx context.Context, minutes int64) ([]domain.NetworkActivity, error) {
	output := make([]domain.NetworkActivity, 0)

	db, err := r.getMongoDB()

	if err != nil {
		return output, err
	}

	cursor, err := db.Collection("activity").Find(ctx, bson.M{"date": bson.M{
		"$gte": primitive.NewDateTimeFromTime(time.Now().Add(-time.Duration(minutes) * time.Minute)),
	}})

	if err != nil {
		return output, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &output)

	return output, err
}

func (r *Repository) GetGeoActivityByMinutes(ctx context.Context, minutes int64) ([]domain.GeoNetworkActivity, error) {
	output := make([]domain.GeoNetworkActivity, 0)

	db, err := r.getMongoDB()

	if err != nil {
		return output, err
	}

	cursor, err := db.Collection("geo_activity").Find(ctx, bson.M{"date": bson.M{
		"$gte": primitive.NewDateTimeFromTime(time.Now().Add(-time.Duration(minutes) * time.Minute)),
	}})

	if err != nil {
		return output, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &output)

	return output, err
}
