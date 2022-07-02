package app

import (
	"context"
	"errors"
	"strings"

	"github.com/lcnssantos/online-activity/internal/domain"
	"github.com/lcnssantos/online-activity/internal/infra/httpclient"
)

type GeoFeatureProperty struct {
	ID       string `json:"id"`
	Oceanic  string `json:"oceanic"`
	LabelLon string `json:"label_lon"`
	LabelLat string `json:"label_lat"`
	Region   string `json:"region"`
	Division string `json:"division"`
}

type GeoFeatureGeometry struct {
	Type        string          `json:"type"`
	Coordinates [][][][]float64 `json:"coordinates"`
}

type GeoFeature struct {
	Type       string             `json:"type"`
	Properties GeoFeatureProperty `json:"properties"`
	Geometry   GeoFeatureGeometry `json:"geometry"`
}

type GeoFIR struct {
	Type     string       `json:"type"`
	Name     string       `json:"name"`
	Features []GeoFeature `json:"features"`
}

type FirService struct {
	geoService   GeoService
	httpClient   httpclient.HttpClient
	responseData *GeoFIR
	endpoint     string
	firsMap      map[string]*domain.FIR
	firsCountry  map[string]*string
}

func NewFirService(geoService GeoService, httpClient httpclient.HttpClient) *FirService {
	return &FirService{
		geoService: geoService,
		httpClient: httpClient,
		endpoint:   "https://map.vatsim.net/livedata/firboundaries.json",
	}
}

func (f *FirService) LoadFirData(ctx context.Context) error {
	if f.responseData == nil {
		var responseData GeoFIR

		err := f.httpClient.Get(ctx, f.endpoint, &responseData)

		if err != nil {
			return err
		}

		for _, geoFIR := range responseData.Features {
			fir := domain.FIR{
				ICAO:   geoFIR.Properties.ID,
				Region: geoFIR.Properties.Region,
				Points: []domain.Point{},
			}

			for _, point := range geoFIR.Geometry.Coordinates[0][0] {
				fir.Points = append(fir.Points, domain.Point{
					Lat: point[1],
					Lon: point[0],
				})
			}

			f.firsMap[fir.ICAO] = &fir
		}

		//TODO - escrever codigo que le os paises
	}

	return nil
}

func (f *FirService) DetectFIR(point domain.Point) (string, error) {
	if f.firsMap == nil {
		return "", errors.New("FIRS NOT LOADED")
	}

	for _, fir := range f.firsMap {
		if f.geoService.IsInside(fir.Points, point) {
			return fir.ICAO, nil
		}
	}

	return "", errors.New("FIR NOT FOUNDED")
}

func (f *FirService) DetectCountryByFIRCode(fir string) string {
	country := f.firsCountry[fir]

	if country == nil {
		return "UNKNOWN"
	}

	return strings.ToUpper(*country)
}

func (f *FirService) DetectCountryByPoint(point domain.Point) (string, error) {
	if f.firsMap == nil {
		return "", errors.New("FIRS NOT LOADED")
	}

	for _, fir := range f.firsMap {
		if f.geoService.IsInside(fir.Points, point) {
			return f.DetectCountryByFIRCode(fir.ICAO), nil
		}
	}

	return "", errors.New("FIR NOT FOUNDED")
}

func (f *FirService) IsInsideFIR(point domain.Point, fir string) (bool, error) {
	if f.firsMap == nil {
		return false, errors.New("FIRS NOT LOADED")
	}

	if f.firsMap[fir] == nil {
		return false, errors.New("FIR NOT FOUNDED")
	}

	return f.geoService.IsInside(f.firsMap[fir].Points, point), nil
}
