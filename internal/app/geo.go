package app

import (
	geo "github.com/kellydunn/golang-geo"
	"github.com/lcnssantos/online-activity/internal/domain"
)

type GeoService struct {
}

func NewGeoService() GeoService {
	return GeoService{}
}

func (g GeoService) IsInside(points []domain.Point, _point domain.Point) bool {
	polygon := geo.NewPolygon(nil)

	for _, point := range points {
		polygon.Add(geo.NewPoint(point.Lat, point.Lon))
	}

	return polygon.Contains(geo.NewPoint(_point.Lat, _point.Lon))
}
