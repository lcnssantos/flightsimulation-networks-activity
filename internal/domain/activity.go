package domain

import (
	"time"
)

type Activity struct {
	Pilot int64 `json:"pilot"`
	ATC   int64 `json:"atc"`
}

type GeoActivity map[string]Activity

type NetworkActivity struct {
	//ID     primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Date   time.Time `json:"date"`
	IVAO   Activity  `json:"ivao"`
	VATSIM Activity  `json:"vatsim"`
	POSCON Activity  `json:"poscon"`
}

type GeoNetworkActivity struct {
	//ID     primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Date   time.Time   `json:"date"`
	IVAO   GeoActivity `json:"ivao"`
	VATSIM GeoActivity `json:"vatsim"`
	POSCON GeoActivity `json:"poscon"`
}
