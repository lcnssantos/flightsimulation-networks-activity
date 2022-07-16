package ivao

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/lcnssantos/online-activity/internal/app"
	"github.com/lcnssantos/online-activity/internal/configuration"
	"github.com/lcnssantos/online-activity/internal/domain"
	"github.com/lcnssantos/online-activity/internal/infra/httpclient"
)

const whazzup = "https://api.ivao.aero/v2/tracker/whazzup"

type cache struct {
	data       *ivaoData
	expiration time.Time
}

type IVAO struct {
	sync.Mutex
	httpClient httpclient.HttpClient
	firService app.FirService
	data       *cache
}

func NewIVAO(httpClient httpclient.HttpClient, firService app.FirService) *IVAO {
	return &IVAO{httpClient: httpClient, firService: firService}
}

func (i *IVAO) loadData(ctx context.Context) (*ivaoData, error) {
	i.Lock()
	defer i.Unlock()

	if i.data != nil && time.Now().Before(i.data.expiration) {
		return i.data.data, nil
	}

	var ivaoData ivaoData

	log.Debug().Msg("Loading IVAO Data")

	err := i.httpClient.Get(ctx, whazzup, &ivaoData)

	if err != nil {
		log.Error().Err(err).Msg("Loading IVAO Data")
		return nil, err
	}

	i.data = &cache{
		data:       &ivaoData,
		expiration: time.Now().Add(time.Minute * time.Duration(configuration.Environment.GetCacheTime())),
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

	wg := sync.WaitGroup{}

	wg.Add(len(data.Clients.Pilots) + len(data.Clients.ATCs))

	for _, pilot := range data.Clients.Pilots {
		go func(pilot ivaoFlight) {
			defer wg.Done()

			if pilot.LastTrack == nil {
				count.increment("UNKNOWN", "pilot")
				return
			}

			country, err := i.firService.DetectCountryByPoint(domain.Point{
				Lat: pilot.LastTrack.Latitude,
				Lon: pilot.LastTrack.Longitude,
			})

			if err != nil {
				count.increment("UNKNOWN", "pilot")
				return
			}

			count.increment(country, "pilot")
		}(pilot)
	}

	for _, atc := range data.Clients.ATCs {
		go func(atc ivaoATC) {
			defer wg.Done()

			if atc.LastTrack == nil {
				count.increment("UNKNOWN", "atc")
				return
			}

			country, err := i.firService.DetectCountryByPoint(domain.Point{
				Lat: atc.LastTrack.Latitude,
				Lon: atc.LastTrack.Longitude,
			})

			if err != nil {
				count.increment("UNKNOWN", "atc")
				return
			}

			count.increment(country, "atc")
		}(atc)
	}

	wg.Wait()

	countData := count.get()

	output := make(domain.GeoActivity)

	for country, activity := range countData {
		output[country] = *activity
	}

	return &output, nil
}
