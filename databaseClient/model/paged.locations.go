package model

type PagedLocation struct {
	MainLocationsCount           *int        `json:"mainLocationsCount"`
	MainLocationsRadiusEnd       *float64    `json:"mainLocationsRadiusEnd"`
	AdditionalLocationsCount     *int        `json:"additionalLocationsCount"`
	AdditionalLocationsRadiusEnd *float64    `json:"additionalLocationsRadiusEnd"`
	MainLocations                *[]Location `json:"mainLocations"`
	AdditionalLocations          *[]Location `json:"additionalLocations"`
}

type Location struct {
	Coordinates  *Coordinates `json:"coordinates"`
	Name         *string      `json:"name"`
	LocationType *string      `json:"type"`
	Distance     *float64     `json:"distance"`
	Time         *int64       `json:"time"`
}
