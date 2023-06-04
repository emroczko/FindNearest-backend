package locations

import (
	"databaseClient/model"
	"github.com/sirupsen/logrus"
	"math"
)

type Service interface {
	GetLocationsByDistance(input *model.LocationByDistanceRequest) (*model.PagedLocation, error)
	GetLocationsByTime(input *model.LocationByTimeRequest) (*model.PagedLocation, error)
}

type service struct {
	repository Repository
}

func NewServiceResult(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) GetLocationsByDistance(request *model.LocationByDistanceRequest) (*model.PagedLocation, error) {

	mainLocationsQueryData := convertQueryParameters(*request.MainLocation)
	var pagedResults model.PagedLocation
	var mainLocationsIds []int64

	for mainLocationsQueryData.RadiusEnd <= *request.MainLocation.RadiusEnd {
		mainLocationsQueryData.Longitude = *request.Longitude
		mainLocationsQueryData.Latitude = *request.Latitude
		mainLocationsResult, mainLocationsErr := s.repository.GetMainLocations(mainLocationsQueryData)

		if mainLocationsErr != nil {
			return nil, mainLocationsErr
		}

		pagedResults.MainLocationsRadiusEnd = &mainLocationsQueryData.RadiusEnd
		var mainLocations []model.Location

		for _, loc := range *mainLocationsResult {
			var locationType string
			if loc.Amenity != nil {
				locationType = *loc.Amenity
			} else if loc.Shop != nil {
				locationType = *loc.Shop
			} else if loc.Water != nil {
				locationType = *loc.Water
			} else if loc.Landuse != nil {
				locationType = *loc.Landuse
			} else {
				locationType = ""
			}

			location := model.Location{
				Coordinates:  loc.Coordinates,
				Name:         loc.Name,
				LocationType: &locationType,
			}
			mainLocations = append(mainLocations, location)
			mainLocationsIds = append(mainLocationsIds, *loc.Osm_Id)
		}

		logrus.Info("End radius: ", mainLocationsQueryData.RadiusEnd)
		logrus.Info("Main locations length: ", len(mainLocations))
		pagedResults.MainLocations = &mainLocations

		if len(mainLocations) != 0 || *request.MainLocation.RadiusEnd == mainLocationsQueryData.RadiusEnd {
			break
		}

		if *request.MainLocation.RadiusEnd-mainLocationsQueryData.RadiusEnd > 7500 {
			mainLocationsQueryData.RadiusEnd += 7500
		} else {
			mainLocationsQueryData.RadiusEnd += *request.MainLocation.RadiusEnd - mainLocationsQueryData.RadiusEnd
		}
	}

	if len(*pagedResults.MainLocations) == 0 {
		pagedResults.MainLocations = &[]model.Location{}
	}

	if request.AdditionalLocation != nil {

		if request.AdditionalLocation.Type == request.MainLocation.Type {
			return &pagedResults, nil
		}

		additionalLocationsQueryData := convertQueryParameters(*request.AdditionalLocation)

		for additionalLocationsQueryData.RadiusEnd <= *request.AdditionalLocation.RadiusEnd {

			pagedResults.AdditionalLocationsRadiusEnd = &additionalLocationsQueryData.RadiusEnd
			var additionalLocations []model.Location

			additionalLocationsResult, additionalLocationsErr := s.repository.GetAdditionalLocations(additionalLocationsQueryData, &mainLocationsIds)

			if additionalLocationsErr != nil {
				return nil, additionalLocationsErr
			}

			for _, loc := range *additionalLocationsResult {
				var locationType string
				if loc.Amenity != nil {
					locationType = *loc.Amenity
				} else if loc.Shop != nil {
					locationType = *loc.Shop
				} else if loc.Water != nil {
					locationType = *loc.Water
				} else if loc.Landuse != nil {
					locationType = *loc.Landuse
				} else {
					locationType = ""
				}

				location := model.Location{
					Coordinates:  loc.Coordinates,
					Name:         loc.Name,
					LocationType: &locationType,
				}
				additionalLocations = append(additionalLocations, location)
			}

			logrus.Info("Additional locations end radius: ", mainLocationsQueryData.RadiusEnd)
			logrus.Info("Additional locations count: ", len(additionalLocations))

			pagedResults.AdditionalLocations = &additionalLocations

			if len(additionalLocations) != 0 || *request.AdditionalLocation.RadiusEnd == additionalLocationsQueryData.RadiusEnd {
				break
			}

			if *request.AdditionalLocation.RadiusEnd-additionalLocationsQueryData.RadiusEnd > 7500 {
				additionalLocationsQueryData.RadiusEnd += 7500
			} else {
				additionalLocationsQueryData.RadiusEnd += *request.AdditionalLocation.RadiusEnd - additionalLocationsQueryData.RadiusEnd
			}
		}

		if len(*pagedResults.AdditionalLocations) == 0 {
			pagedResults.AdditionalLocations = &[]model.Location{}
		}
	} else {
		additionalLocationsRadiusEnd := 0.0
		pagedResults.AdditionalLocationsRadiusEnd = &additionalLocationsRadiusEnd
		pagedResults.AdditionalLocations = &[]model.Location{}
	}

	return &pagedResults, nil
}

func (s *service) GetLocationsByTime(request *model.LocationByTimeRequest) (*model.PagedLocation, error) {
	return nil, nil
}

func convertQueryParameters(request model.LocationByDistanceRequestDetails) *model.LocationQuery {
	var optimizedRequest model.LocationQuery

	optimizedRequest.Type = *request.Type

	if *request.RadiusEnd > 7500 && *request.RadiusStart == 0 {
		optimizedRequest.RadiusEnd = 2000 * math.Log10(*request.RadiusEnd)
		optimizedRequest.RadiusStart = 0
		logrus.Info("Optimizing end radius, start radius == 0")
	} else if *request.RadiusEnd > 7500 && *request.RadiusStart > 7500 {
		if *request.RadiusEnd-*request.RadiusStart < 7500 {
			optimizedRequest.RadiusEnd = *request.RadiusEnd
			optimizedRequest.RadiusStart = *request.RadiusStart
			logrus.Info("Original radius bigger than 7500")
		} else {
			optimizedRequest.RadiusEnd = 2000 * math.Log10(*request.RadiusEnd)
			optimizedRequest.RadiusStart = *request.RadiusStart
			logrus.Info("Optimizing end radius, start radius != 0")
		}
	} else {
		optimizedRequest.RadiusEnd = *request.RadiusEnd
		optimizedRequest.RadiusStart = *request.RadiusStart
		logrus.Info("Original radius")
	}

	return &optimizedRequest
}
