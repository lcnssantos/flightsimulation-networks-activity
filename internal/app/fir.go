package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lcnssantos/online-activity/internal/domain"
	"github.com/lcnssantos/online-activity/internal/infra/httpclient"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
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

type FirCountry struct {
	ICAO    string `json:"ICAO"`
	Country string `json:"Country"`
}

func NewFirService(geoService GeoService, httpClient httpclient.HttpClient) FirService {
	return FirService{
		geoService:  geoService,
		httpClient:  httpClient,
		endpoint:    "https://map.vatsim.net/livedata/firboundaries.json",
		firsCountry: map[string]*string{},
		firsMap:     map[string]*domain.FIR{},
	}
}

func (f *FirService) loadFirCountryData() error {
	log.Info().Msg("Start FIR country loading...")

	path, err := os.Getwd()

	if err != nil {
		return err
	}

	file := fmt.Sprintf("%s/data/firs_countries.json", path)

	data, err := os.ReadFile(file)

	if err != nil {
		return err
	}

	countries := make([]FirCountry, 0)

	err = json.Unmarshal(data, &countries)

	if err != nil {
		return err
	}

	for _, country := range countries {
		countryName := country.Country
		f.firsCountry[country.ICAO] = &countryName
	}

	return nil
}

func (f *FirService) LoadFirData(ctx context.Context) error {
	if f.responseData == nil {
		log.Info().Msg("Starting FIR data loading...")

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

		err = f.loadFirCountryData()

		if err != nil {
			return err
		}
	}

	return nil
}

func (f *FirService) DetectFIR(point domain.Point) (string, error) {
	log.Debug().Interface("point", point).Msg("Detect FIR from point")

	if f.firsMap == nil {
		log.Error().Msg("FIRs not loaded")
		return "", errors.New("FIRS NOT LOADED")
	}

	for _, fir := range f.firsMap {
		if f.geoService.IsInside(fir.Points, point) {
			return fir.ICAO, nil
		}
	}

	log.Warn().Interface("point", point).Msg("FIR not founded")

	return "", errors.New("FIR NOT FOUNDED")
}

func (f *FirService) DetectCountryByFIRCode(fir string) string {
	log.Debug().Interface("fir", fir).Msg("Get country by fir code")

	country := f.firsCountry[fir]

	if country == nil {
		log.Warn().Interface("fir", fir).Msg("Country not founded")
		return "UNKNOWN"
	}

	return strings.ToUpper(*country)
}

func (f *FirService) DetectCountryByPoint(point domain.Point) (string, error) {
	log.Debug().Interface("point", point).Msg("Get country by point")

	if f.firsMap == nil {
		log.Error().Msg("FIRs not loaded")
		return "", errors.New("FIRS NOT LOADED")
	}

	for _, fir := range f.firsMap {
		if f.geoService.IsInside(fir.Points, point) {
			return f.DetectCountryByFIRCode(fir.ICAO), nil
		}
	}

	log.Warn().Interface("point", point).Msg("Fir not founded")

	return "", errors.New("FIR NOT FOUNDED")
}

func (f *FirService) IsInsideFIR(point domain.Point, fir string) (bool, error) {
	if f.firsMap == nil {
		log.Error().Interface("point", point).Interface("fir", fir).Msg("FIRs not loaded")
		return false, errors.New("FIRS NOT LOADED")
	}

	if f.firsMap[fir] == nil {
		log.Error().Interface("point", point).Interface("fir", fir).Msg("FIR Not founded")
		return false, errors.New("FIR NOT FOUNDED")
	}

	return f.geoService.IsInside(f.firsMap[fir].Points, point), nil
}
