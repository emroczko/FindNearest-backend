package model

import (
	"databaseClient/util"
)

type Location struct {
	Amenity          *string
	Name             *string
	Shop             *string
	Sport            *string
	Public_Transport *string
	Tags             *string
	Water            *string
	Landuse          *string
	Coordinates      *util.Coordinates
}
