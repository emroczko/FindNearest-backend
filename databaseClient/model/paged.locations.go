package model

type PagedLocation struct {
	RadiusEnd float64    `json:"radiusEnd"`
	Locations []Location `json:"locations"`
}

type Location struct {
	Coordinates  Coordinates `json:"coordinates"`
	Name         string      `json:"name"`
	LocationType string      `json:"type"`
}
