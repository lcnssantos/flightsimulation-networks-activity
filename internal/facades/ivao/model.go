package ivao

type ivaoLastTrack struct {
	latitude  float64 `json:"latitude"`
	longitude float64 `json:"longitude"`
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
