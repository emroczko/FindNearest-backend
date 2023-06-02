package locations

import (
	"databaseClient/model"
	"github.com/sirupsen/logrus"
	"math"
)

type Service interface {
	ResultLocationsService(input *model.LocationByDistanceRequest) (*model.PagedLocation, error)
}

type service struct {
	repository Repository
}

func NewServiceResult(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) ResultLocationsService(request *model.LocationByDistanceRequest) (*model.PagedLocation, error) {

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
			location := model.Location{
				Coordinates:  loc.Coordinates,
				Name:         loc.Name,
				LocationType: loc.Amenity,
			}
			mainLocations = append(mainLocations, location)
			mainLocationsIds = append(mainLocationsIds, *loc.Osm_Id)
		}

		pagedResults.MainLocations = &mainLocations

		if len(mainLocations) != 0 {
			break
		}

		mainLocationsQueryData.RadiusEnd += 7500
		if *request.MainLocation.RadiusEnd-mainLocationsQueryData.RadiusEnd > 7500 {
			mainLocationsQueryData.RadiusEnd += 7500
		} else {
			mainLocationsQueryData.RadiusEnd += *request.MainLocation.RadiusEnd - mainLocationsQueryData.RadiusEnd
		}
	}

	if request.AdditionalLocation != nil {

		if request.AdditionalLocation.Type == request.MainLocation.Type {
			return &pagedResults, nil
		}

		additionalLocationsQueryData := convertQueryParameters(*request.AdditionalLocation)

		for additionalLocationsQueryData.RadiusEnd <= *request.AdditionalLocation.RadiusEnd {
			var additionalLocations []model.Location

			additionalLocationsResult, additionalLocationsErr := s.repository.GetAdditionalLocations(additionalLocationsQueryData, &mainLocationsIds)

			if additionalLocationsErr != nil {
				return nil, additionalLocationsErr
			}

			for _, loc := range *additionalLocationsResult {
				location := model.Location{
					Coordinates:  loc.Coordinates,
					Name:         loc.Name,
					LocationType: loc.Amenity,
				}
				additionalLocations = append(additionalLocations, location)
			}

			pagedResults.AdditionalLocations = &additionalLocations

			if len(additionalLocations) != 0 {
				break
			}

			additionalLocationsQueryData.RadiusEnd += 7500
			if *request.AdditionalLocation.RadiusEnd-additionalLocationsQueryData.RadiusEnd > 7500 {
				additionalLocationsQueryData.RadiusEnd += 7500
			} else {
				additionalLocationsQueryData.RadiusEnd += *request.AdditionalLocation.RadiusEnd - additionalLocationsQueryData.RadiusEnd
			}
		}
	}

	return &pagedResults, nil
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
