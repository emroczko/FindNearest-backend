package model

type LocationByDistanceRequest struct {
	Latitude           *float64                          `json:"latitude" binding:"required"`
	Longitude          *float64                          `json:"longitude" binding:"required"`
	MainLocation       *LocationByDistanceRequestDetails `json:"mainLocation" binding:"required"`
	AdditionalLocation *LocationByDistanceRequestDetails `json:"additionalLocation" binding:"omitempty"`
}

type LocationByDistanceRequestDetails struct {
	Type        *string  `json:"type" binding:"required"`
	RadiusStart *float64 `json:"radiusStart" binding:"min=0"`
	RadiusEnd   *float64 `json:"radiusEnd" binding:"required,max=10000"`
}

type LocationByTimeRequest struct {
	Latitude      *float64                          `json:"latitude" binding:"required"`
	Longitude     *float64                          `json:"longitude" binding:"required"`
	MainLocation  *LocationByTimeRequestDetails     `json:"mainLocation" binding:"required"`
	OtherLocation *LocationByDistanceRequestDetails `json:"otherLocation" binding:"omitempty"`
}

type LocationByTimeRequestDetails struct {
	Type      *string  `json:"type" binding:"required"`
	TimeStart *float64 `json:"timeStart" binding:"min=0"`
	TimeEnd   *float64 `json:"timeEnd" binding:"required,max=10000"`
}
