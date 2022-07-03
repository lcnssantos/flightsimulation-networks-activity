package vatsim

type vatsimTransceiverDataTransceiver struct {
	Latitude  *float64 `json:"latDeg"`
	Longitude *float64 `json:"lonDeg"`
}

type vatsimTransceiverData struct {
	Callsign     *string                            `json:"callsign"`
	Transceivers []vatsimTransceiverDataTransceiver `json:"transceivers"`
}

type vatsimPilot struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type vatsimATC struct {
	Callsign string `json:"callsign"`
}

type vatsimData struct {
	Pilots []vatsimPilot `json:"pilots"`
	Atc    []vatsimATC   `json:"controllers"`
}
