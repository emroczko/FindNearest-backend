package model

type LocationsTimes struct {
	LocationDetails *PossibleLocationDetails `json:"locationsDetails" binding:"required"`
	Distance        *float64                 `json:"distance" binding:"required"`
	Time            *int64                   `json:"time" binding:"required"`
}
