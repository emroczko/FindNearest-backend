package model

type RouteToLocation struct {
	SourceCoordinates        *Coordinates               `json:"sourceCoordinates"`
	PossibleLocationsDetails *[]PossibleLocationDetails `json:"locationsDetails"`
	MeanOfTransport          *string                    `json:"meanOfTransport"`
}

type PossibleLocationDetails struct {
	Id          int64        `json:"id"`
	Coordinates *Coordinates `json:"coordinates"`
}
