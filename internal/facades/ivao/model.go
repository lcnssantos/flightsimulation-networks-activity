package ivao

type ivaoLastTrack struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ivaoFlight struct {
	LastTrack *ivaoLastTrack `json:"lastTrack"`
}

type ivaoATC struct {
	LastTrack *ivaoLastTrack `json:"lastTrack"`
}

type ivaoClients struct {
	Pilots []ivaoFlight `json:"pilots"`
	ATCs   []ivaoATC    `json:"atcs"`
}

type ivaoData struct {
	Clients ivaoClients `json:"clients"`
}
