package model

type LocationQuery struct {
	Latitude            float64
	Longitude           float64
	Type                string
	RadiusStart         float64
	RadiusEnd           float64
	OriginalRadiusStart float64
	OriginalRadiusEnd   float64
}
