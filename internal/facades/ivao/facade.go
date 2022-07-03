package ivao

import (
	"context"
	"github.com/rs/zerolog/log"

	"github.com/lcnssantos/online-activity/internal/app"
	"github.com/lcnssantos/online-activity/internal/domain"
	"github.com/lcnssantos/online-activity/internal/infra/httpclient"
)

const whazzup = "https://api.ivao.aero/v2/tracker/whazzup"

type IVAO struct {
	httpClient httpclient.HttpClient
	firService app.FirService
}

func NewIVAO(httpClient httpclient.HttpClient, firService app.FirService) *IVAO {
	return &IVAO{httpClient: httpClient, firService: firService}
}

func (i *IVAO) loadData(ctx context.Context) (*ivaoData, error) {
	var ivaoData ivaoData

	log.Debug().Msg("Loading IVAO Data")

	err := i.httpClient.Get(ctx, whazzup, &ivaoData)

	if err != nil {
		log.Error().Err(err).Msg("Loading IVAO Data")
		return nil, err
	}

	return &ivaoData, nil
}

func (i *IVAO) GetActivity(ctx context.Context) (*domain.Activity, error) {
	log.Debug().Msg("Get IVAO Activity")

	data, err := i.loadData(ctx)

	if err != nil {
		log.Error().Err(err).Msg("Get IVAO Activity")
		return nil, err
	}

	return &domain.Activity{
		Pilot: int64(len(data.Clients.Pilots)),
		ATC:   int64(len(data.Clients.ATCs)),
	}, nil
}

func (i *IVAO) GetBrazilActivity(ctx context.Context) (*domain.Activity, error) {
	log.Debug().Msg("Get IVAO Brazil Activity")

	data, err := i.loadData(ctx)

	if err != nil {
		log.Error().Err(err).Msg("Get IVAO Brazil Activity")
		return nil, err
	}

	activity := domain.Activity{
		Pilot: 0,
		ATC:   0,
	}

	firs := []string{"SBBS", "SBCW", "SBRE", "SBAZ", "SBAO"}

	for _, pilot := range data.Clients.Pilots {
		if pilot.LastTrack != nil {
			for _, fir := range firs {
				isInsideFir, err := i.firService.IsInsideFIR(domain.Point{
					Lat: pilot.LastTrack.Latitude,
					Lon: pilot.LastTrack.Longitude,
				}, fir)

				if err != nil {
					return nil, err
				}

				if isInsideFir {
					activity.Pilot++
					break
				}

			}
		}
	}

	for _, atc := range data.Clients.ATCs {
		if atc.LastTrack == nil {
			break
		}

		for _, fir := range firs {
			isInsideFir, err := i.firService.IsInsideFIR(domain.Point{
				Lat: atc.LastTrack.Latitude,
				Lon: atc.LastTrack.Longitude,
			}, fir)

			if err != nil {
				return nil, err
			}

			if isInsideFir {
				activity.ATC++
				break
			}

		}
	}

	return &activity, nil
}

func (i *IVAO) GetGeoActivity(ctx context.Context) (*domain.GeoActivity, error) {
	log.Debug().Msg("Get IVAO Geo Activity")

	data, err := i.loadData(ctx)

	if err != nil {
		log.Error().Err(err).Msg("Get IVAO Geo Activity")
		return nil, err
	}

	count := newCount()

	for _, pilot := range data.Clients.Pilots {
		if pilot.LastTrack == nil {
			count.increment("UNKNOWN", "pilot")
			break
		}

		country, err := i.firService.DetectCountryByPoint(domain.Point{
			Lat: pilot.LastTrack.Latitude,
			Lon: pilot.LastTrack.Longitude,
		})

		if err != nil {
			count.increment("UNKNOWN", "pilot")
		} else {
			count.increment(country, "pilot")
		}
	}

	for _, atc := range data.Clients.ATCs {
		if atc.LastTrack == nil {
			count.increment("UNKNOWN", "atc")
			break
		}

		country, err := i.firService.DetectCountryByPoint(domain.Point{
			Lat: atc.LastTrack.Latitude,
			Lon: atc.LastTrack.Longitude,
		})

		if err != nil {
			count.increment("UNKNOWN", "atc")
		} else {
			count.increment(country, "atc")
		}
	}

	countData := count.get()

	output := make(domain.GeoActivity)

	for country, activity := range countData {
		output[country] = *activity
	}

	return &output, nil
}
