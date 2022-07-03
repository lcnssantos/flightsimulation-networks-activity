package vatsim

import (
	"context"
	"github.com/rs/zerolog/log"

	"github.com/lcnssantos/online-activity/internal/app"
	"github.com/lcnssantos/online-activity/internal/domain"
	"github.com/lcnssantos/online-activity/internal/infra/httpclient"
)

const whazzup = "https://data.vatsim.net/v3/vatsim-data.json"
const transceiverHost = "https://data.vatsim.net/v3/transceivers-data.json"

type VATSIM struct {
	httpClient      httpclient.HttpClient
	firService      app.FirService
	transceiverData map[string]*domain.Point
}

func NewVatsim(httpCLient httpclient.HttpClient, firService app.FirService) *VATSIM {
	return &VATSIM{httpClient: httpCLient, firService: firService, transceiverData: map[string]*domain.Point{}}
}

func (v *VATSIM) loadData(ctx context.Context) (*vatsimData, error) {
	var transceiverData []vatsimTransceiverData

	log.Debug().Msg("Loading VATSIM Data")

	err := v.httpClient.Get(ctx, transceiverHost, &transceiverData)

	if err != nil {
		log.Error().Err(err).Msg("Loading VATSIM Data")
		return nil, err
	}

	for _, transceiver := range transceiverData {
		if len(transceiver.Transceivers) > 0 {
			v.transceiverData[*transceiver.Callsign] = &domain.Point{
				Lat: *transceiver.Transceivers[0].Latitude,
				Lon: *transceiver.Transceivers[0].Longitude,
			}

		}
	}

	var data vatsimData

	err = v.httpClient.Get(ctx, whazzup, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (v *VATSIM) GetActivity(ctx context.Context) (*domain.Activity, error) {
	log.Debug().Msg("Get VATSIM Activity")

	data, err := v.loadData(ctx)

	if err != nil {
		log.Error().Err(err).Msg("Get VATSIM Activity")
		return nil, err
	}

	return &domain.Activity{
		Pilot: int64(len(data.Pilots)),
		ATC:   int64(len(data.Atc)),
	}, nil
}

func (v *VATSIM) GetBrazilActivity(ctx context.Context) (*domain.Activity, error) {
	log.Debug().Msg("Get VATSIM Brazil Activity")

	data, err := v.loadData(ctx)

	if err != nil {
		log.Error().Err(err).Msg("Get VATSIM Brazil Activity")
		return nil, err
	}

	activity := domain.Activity{
		Pilot: 0,
		ATC:   0,
	}

	firs := []string{"SBBS", "SBCW", "SBRE", "SBAZ", "SBAO"}

	for _, pilot := range data.Pilots {
		if pilot.Latitude != nil && pilot.Longitude != nil {
			for _, fir := range firs {
				isInsideFir, err := v.firService.IsInsideFIR(domain.Point{
					Lat: *pilot.Latitude,
					Lon: *pilot.Longitude,
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

	for _, atc := range data.Atc {
		if atc.Callsign[0] == 'S' && (atc.Callsign[1] == 'B' || atc.Callsign[1] == 'D') {
			activity.ATC++
		}
	}

	return &activity, nil
}

func (v *VATSIM) GetGeoActivity(ctx context.Context) (*domain.GeoActivity, error) {
	log.Debug().Msg("Get VATSIM Geo Activity")

	data, err := v.loadData(ctx)

	if err != nil {
		log.Error().Err(err).Msg("Get VATSIM Geo Activity")
		return nil, err
	}

	count := newCount()

	for _, pilot := range data.Pilots {
		if pilot.Latitude == nil || pilot.Longitude == nil {
			count.increment("UNKNOWN", "pilot")
			break
		}

		country, err := v.firService.DetectCountryByPoint(domain.Point{
			Lat: *pilot.Latitude,
			Lon: *pilot.Longitude,
		})

		if err != nil {
			count.increment("UNKNOWN", "pilot")
		} else {
			count.increment(country, "pilot")
		}
	}

	for _, atc := range data.Atc {
		point := v.transceiverData[atc.Callsign]

		if point == nil {
			count.increment("UNKNOWN", "atc")
			break
		}

		country, err := v.firService.DetectCountryByPoint(*point)

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
