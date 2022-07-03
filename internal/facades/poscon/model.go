package poscon

import "time"

type posconFlightPosition struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type posconFlight struct {
	Position *posconFlightPosition `json:"position"`
}

type posconATC struct {
	CenterPoint *[]float64 `json:"centerPoint"`
	Fir         string     `json:"fir"`
}

type posconData struct {
	TotalPilots int64          `json:"totalPilots"`
	TotalAtc    int64          `json:"totalAtc"`
	LastUpdated time.Time      `json:"lastUpdated"`
	Flights     []posconFlight `json:"flights"`
	Atc         []posconATC    `json:"atc"`
}
