package locations

type LocationRequest struct {
	Latitude     float64      `json:"latitude" binding:"required"`
	Longitude    float64      `json:"longitude" binding:"required"`
	Type         string       `json:"type" binding:"required"`
	Distance     float64      `json:"distance" binding:"required"`
	DistanceUnit DistanceUnit `json:"unit" binding:"required"`
}

type DistanceUnit string

const (
	Meters     DistanceUnit = "meters"
	Kilometers DistanceUnit = "kilometers"
)
