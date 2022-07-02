package main

import (
	"context"
	"github.com/lcnssantos/online-activity/internal/app"
	"github.com/lcnssantos/online-activity/internal/facades/ivao"
	"github.com/lcnssantos/online-activity/internal/facades/poscon"
	"github.com/lcnssantos/online-activity/internal/facades/vatsim"
	"github.com/lcnssantos/online-activity/internal/infra/database"
	"github.com/lcnssantos/online-activity/internal/infra/httpclient"
	"github.com/lcnssantos/online-activity/internal/infra/httpserver"
)

func main() {
	httpClient := httpclient.NewHttpClient()
	geoService := app.NewGeoService()
	firService := app.NewFirService(geoService, httpClient)

	err := firService.LoadFirData(context.Background())

	if err != nil {
		panic(err)
	}

	ivaoFacade := ivao.NewIVAO(httpClient, firService)
	vatsimFacade := vatsim.NewVatsim(httpClient, firService)
	posconFacade := poscon.NewPoscon(httpClient, firService)

	mongo := database.NewMongoDriver("mongodb://localhost:27017", "tracking")

	repository := database.NewRepository(mongo)

	appService := app.NewAppService(ivaoFacade, posconFacade, vatsimFacade, &repository)

	controller := httpserver.NewController(appService)

	server := httpserver.NewServer(8080, controller)

	server.Listen()
}
