package main

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/lcnssantos/online-activity/internal/app"
	"github.com/lcnssantos/online-activity/internal/configuration"
	"github.com/lcnssantos/online-activity/internal/facades/ivao"
	"github.com/lcnssantos/online-activity/internal/facades/poscon"
	"github.com/lcnssantos/online-activity/internal/facades/vatsim"
	"github.com/lcnssantos/online-activity/internal/infra/database"
	"github.com/lcnssantos/online-activity/internal/infra/httpclient"
	"github.com/lcnssantos/online-activity/internal/infra/httpserver"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Panic().Msg("Error loading .env file")
	// }

	configuration.LoadEnv()

	err := configuration.Environment.Validate()

	if err != nil {
		log.Panic().Err(err).Msg("Error to load environment variables")
	}

	httpClient := httpclient.NewHttpClient()
	geoService := app.NewGeoService()
	firService := app.NewFirService(geoService, httpClient)

	err = firService.LoadFirData(context.Background())

	if err != nil {
		log.Panic().Err(err).Msg("Error to load fir data")
	}

	ivaoFacade := ivao.NewIVAO(httpClient, firService)
	vatsimFacade := vatsim.NewVatsim(httpClient, firService)
	posconFacade := poscon.NewPoscon(httpClient, firService)

	mongo := database.NewMongoDriver(configuration.Environment.MongoURL, "tracking")

	repository := database.NewRepository(mongo)

	appService := app.NewAppService(ivaoFacade, posconFacade, vatsimFacade, &repository)

	controller := httpserver.NewController(appService)

	server := httpserver.NewServer(configuration.Environment.GetPort(), controller)

	server.Listen()
}
