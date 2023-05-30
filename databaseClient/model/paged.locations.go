package model

type PagedLocation struct {
	PagesCount int
	PageNumber int
	Locations  []Location
}

type Location struct {
	Coordinates Coordinates
	Name        string
	Amenity     string
}
