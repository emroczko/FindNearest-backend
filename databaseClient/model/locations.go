package model

type LocationEntity struct {
	Coordinates      *Coordinates
	Name             *string
	Amenity          *string
	Shop             *string
	Sport            *string
	Public_Transport *string
	Tags             *string
	Water            *string
	Landuse          *string
}
