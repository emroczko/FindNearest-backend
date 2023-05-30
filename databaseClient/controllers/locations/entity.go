package locations

type LocationRequest struct {
	Latitude    float64 `form:"latitude" binding:"required"`
	Longitude   float64 `form:"longitude" binding:"required"`
	Type        string  `form:"type" binding:"required"`
	RadiusStart float64 `form:"radiusStart" binding:"min=0"`
	RadiusEnd   float64 `form:"radiusEnd" binding:"required,max=5000"`
}
