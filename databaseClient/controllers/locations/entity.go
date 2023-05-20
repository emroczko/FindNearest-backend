package locations

type LocationRequest struct {
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	RadiusStart float64 `json:"radiusStart"`
	RadiusEnd   float64 `json:"radiusEnd" binding:"required"`
}
