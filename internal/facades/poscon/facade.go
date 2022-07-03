package poscon

import (
	"context"
	"github.com/rs/zerolog/log"

	"github.com/lcnssantos/online-activity/internal/app"
	"github.com/lcnssantos/online-activity/internal/domain"
	"github.com/lcnssantos/online-activity/internal/infra/httpclient"
)

const whazzup = "https://hqapi.poscon.net/online.json"

type Poscon struct {
	httpClient httpclient.HttpClient
	firService app.FirService
}

func NewPoscon(httpCLient httpclient.HttpClient, firService app.FirService) *Poscon {
	return &Poscon{httpClient: httpCLient, firService: firService}
}

func (p *Poscon) loadData(ctx context.Context) (*posconData, error) {
	var posconData posconData

	log.Debug().Msg("Loading POSCON Data")

	err := p.httpClient.Get(ctx, whazzup, &posconData)

	if err != nil {
		log.Error().Err(err).Msg("Loading POSCON Data")
		return nil, err
	}

	return &posconData, nil
}

func (p *Poscon) GetActivity(ctx context.Context) (*domain.Activity, error) {
	log.Debug().Msg("Get POSCON Activity")

	data, err := p.loadData(ctx)

	if err != nil {
		log.Error().Err(err).Msg("Get POSCON Activity")
		return nil, err
	}

	return &domain.Activity{
		Pilot: data.TotalPilots,
		ATC:   int64(data.TotalAtc),
	}, nil
}

func (p *Poscon) GetBrazilActivity(ctx context.Context) (*domain.Activity, error) {
	log.Debug().Msg("Get POSCON Brazil Activity")

	data, err := p.loadData(ctx)

	if err != nil {
		log.Error().Err(err).Msg("Get POSCON Brazil Activity")
		return nil, err
	}

	activity := domain.Activity{
		Pilot: 0,
		ATC:   0,
	}

	firs := []string{"SBBS", "SBCW", "SBRE", "SBAZ", "SBAO"}

	for _, pilot := range data.Flights {
		if pilot.Position != nil {
			for _, fir := range firs {
				isInsideFir, err := p.firService.IsInsideFIR(domain.Point{
					Lat: pilot.Position.Lat,
					Lon: pilot.Position.Long,
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
		for _, fir := range firs {
			if fir == atc.Fir {
				activity.ATC++
				break
			}
		}
	}

	return &activity, nil
}

func (p *Poscon) GetGeoActivity(ctx context.Context) (*domain.GeoActivity, error) {
	log.Debug().Msg("Get POSCON Geo Activity")

	data, err := p.loadData(ctx)

	if err != nil {
		log.Error().Err(err).Msg("Get POSCON Geo Activity")
		return nil, err
	}

	count := newCount()

	for _, pilot := range data.Flights {
		if pilot.Position == nil {
			count.increment("UNKNOWN", "pilot")
			break
		}

		country, err := p.firService.DetectCountryByPoint(domain.Point{
			Lat: pilot.Position.Lat,
			Lon: pilot.Position.Long,
		})

		if err != nil {
			count.increment("UNKNOWN", "pilot")
		} else {
			count.increment(country, "pilot")
		}
	}

	for _, atc := range data.Atc {
		if atc.CenterPoint == nil || len(*atc.CenterPoint) == 0 {
			count.increment("UNKNOWN", "atc")
			break
		}

		country := p.firService.DetectCountryByFIRCode(atc.Fir)

		count.increment(country, "atc")
	}

	countData := count.get()

	output := make(domain.GeoActivity)

	for country, activity := range countData {
		output[country] = *activity
	}

	return &output, nil
}
