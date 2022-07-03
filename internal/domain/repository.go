package domain

import "context"

type Repository interface {
	SaveBrazilActivity(ctx context.Context, activity NetworkActivity) error
	SaveActivity(ctx context.Context, activity NetworkActivity) error
	SaveGeoActivity(ctx context.Context, activity GeoNetworkActivity) error
	GetBrazilActivityByMinutes(ctx context.Context, minutes int64) ([]NetworkActivity, error)
	GetActivityByMinutes(ctx context.Context, minutes int64) ([]NetworkActivity, error)
	GetGeoActivityByMinutes(ctx context.Context, minutes int64) ([]GeoNetworkActivity, error)
}
