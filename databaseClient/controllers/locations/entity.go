package locations

type LocationRequest struct {
	Latitude     float64      `json:"latitude" binding:"required"`
	Longitude    float64      `json:"longitude" binding:"required"`
	Type         LocationType `json:"type" binding:"required"`
	Distance     float64      `json:"distance" binding:"required"`
	DistanceUnit DistanceUnit `json:"unit" binding:"required"`
}

type LocationType string

const (
	Shop   LocationType = "shop"
	Church LocationType = "church"
	School LocationType = "school"
)

type DistanceUnit string

const (
	Meters     LocationType = "meters"
	Kilometers LocationType = "kilometers"
)
