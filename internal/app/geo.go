package app

import (
	"github.com/lcnssantos/online-activity/internal/domain"
)

type GeoService struct {
}

func NewGeoService() GeoService {
	return GeoService{}
}

func (g GeoService) IsInside(points []domain.Point, _point domain.Point) bool {
	polygon := make([][]float64, 0)

	for _, p := range points {
		polygon = append(polygon, []float64{p.Lon, p.Lat})
	}

	point := []float64{_point.Lon, _point.Lat}

	result := false

	for i, j := 0, len(polygon)-1; i < len(polygon); i++ {
		if (polygon[i][1] > point[1]) != (polygon[j][1] > point[1]) &&
			point[0] < ((polygon[j][0]-polygon[i][0])*(point[1]-polygon[i][1]))/(polygon[j][1]-polygon[i][1])+polygon[i][0] {
			result = !result
		}
		j = i
	}

	return result
}
