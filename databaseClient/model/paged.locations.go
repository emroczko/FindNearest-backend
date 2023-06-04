package model

type PagedLocation struct {
	MainLocationsRadiusEnd       *float64    `json:"mainLocationsRadiusEnd"`
	AdditionalLocationsRadiusEnd *float64    `json:"additionalLocationsRadiusEnd"`
	MainLocations                *[]Location `json:"mainLocations"`
	AdditionalLocations          *[]Location `json:"additionalLocations"`
}

type Location struct {
	Coordinates  *Coordinates `json:"coordinates"`
	Name         *string      `json:"name"`
	LocationType *string      `json:"type"`
}
