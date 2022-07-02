package domain

import "context"

type Facade interface {
	GetActivity(ctx context.Context) (*Activity, error)
	GetBrazilActivity(ctx context.Context) (*Activity, error)
	GetGeoActivity(ctx context.Context) (*GeoActivity, error)
}
