package locations

import (
	"bytes"
	"databaseClient/config"
	"databaseClient/model"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
	"net/http"
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

	mainLocationsQueryData := optimizeDistanceQueryParameters(request.MainLocation.RadiusStart, request.MainLocation.RadiusEnd)
	var pagedResults model.PagedLocation
	var mainLocationsIds []int64

	for mainLocationsQueryData.RadiusEnd <= *request.MainLocation.RadiusEnd {
		mainLocationsQueryData.Longitude = *request.Longitude
		mainLocationsQueryData.Latitude = *request.Latitude
		mainLocationsQueryData.Type = *request.MainLocation.Type
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

		additionalLocations, additionalLocationsEndRadius, err := s.getAdditionalLocations(request.AdditionalLocation, &mainLocationsIds)

		if err != nil {
			return nil, err
		}

		pagedResults.AdditionalLocations = additionalLocations
		pagedResults.AdditionalLocationsRadiusEnd = additionalLocationsEndRadius

	} else {
		additionalLocationsRadiusEnd := 0.0
		pagedResults.AdditionalLocationsRadiusEnd = &additionalLocationsRadiusEnd
		pagedResults.AdditionalLocations = &[]model.Location{}
	}

	return &pagedResults, nil
}

func (s *service) GetLocationsByTime(request *model.LocationByTimeRequest) (*model.PagedLocation, error) {

	var radiusStart = float64(*request.MainLocation.TimeStart * 600)
	var radiusEnd = float64(*request.MainLocation.TimeEnd * 600)

	mainLocationsQueryData := optimizeDistanceQueryParameters(&radiusStart, &radiusEnd)
	mainLocationsQueryData.Type = *request.MainLocation.Type

	var pagedResults model.PagedLocation
	var mainLocationsIds []int64
	mainLocations := make(map[int64]*model.Location)

	for mainLocationsQueryData.RadiusEnd <= radiusEnd {
		mainLocationsQueryData.Longitude = *request.Longitude
		mainLocationsQueryData.Latitude = *request.Latitude
		mainLocationsQueryData.Type = *request.MainLocation.Type
		mainLocationsResult, mainLocationsErr := s.repository.GetMainLocations(mainLocationsQueryData)

		if mainLocationsErr != nil {
			return nil, mainLocationsErr
		}

		pagedResults.MainLocationsRadiusEnd = &mainLocationsQueryData.RadiusEnd

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
			mainLocations[*loc.Osm_Id] = &location
			mainLocationsIds = append(mainLocationsIds, *loc.Osm_Id)
		}

		logrus.Info("End radius: ", mainLocationsQueryData.RadiusEnd)
		logrus.Info("Main locations length: ", len(mainLocations))

		if len(mainLocations) != 0 || radiusEnd == mainLocationsQueryData.RadiusEnd {
			break
		}

		if radiusEnd-mainLocationsQueryData.RadiusEnd > 7500 {
			mainLocationsQueryData.RadiusEnd += 7500
		} else {
			mainLocationsQueryData.RadiusEnd += radiusEnd - mainLocationsQueryData.RadiusEnd
		}
	}

	if len(mainLocations) == 0 {
		pagedResults.MainLocations = &[]model.Location{}
	} else {
		locationsTimes, timesErr := getTimeToLocations(request.Latitude, request.Longitude, &mainLocations)

		if timesErr != nil {
			return nil, timesErr
		}

		var checkedMainLocations []model.Location

		for _, loc := range *locationsTimes {
			fmt.Println(*loc.LocationDetails.Id)
			if *loc.Time > *request.MainLocation.TimeEnd {
				fmt.Println(mainLocations[*loc.LocationDetails.Id])
				checkedMainLocations = append(checkedMainLocations, model.Location{
					Coordinates:  loc.LocationDetails.Coordinates,
					Name:         mainLocations[*loc.LocationDetails.Id].Name,
					LocationType: mainLocations[*loc.LocationDetails.Id].LocationType,
				})
			}
		}

		pagedResults.MainLocations = &checkedMainLocations

		if request.AdditionalLocation != nil {

			if request.AdditionalLocation.Type == request.MainLocation.Type {
				return &pagedResults, nil
			}

			additionalLocations, additionalLocationsEndRadius, err := s.getAdditionalLocations(request.AdditionalLocation, &mainLocationsIds)

			if err != nil {
				return nil, err
			}

			pagedResults.AdditionalLocations = additionalLocations
			pagedResults.AdditionalLocationsRadiusEnd = additionalLocationsEndRadius

		} else {
			additionalLocationsRadiusEnd := 0.0
			pagedResults.AdditionalLocationsRadiusEnd = &additionalLocationsRadiusEnd
			pagedResults.AdditionalLocations = &[]model.Location{}
		}

	}
	return &pagedResults, nil
}

func getTimeToLocations(sourceLatitude *float64, sourceLongitude *float64, locationsWithId *map[int64]*model.Location) (*[]model.LocationsTimes, error) {

	var possibleLocationsCoordinates []model.PossibleLocationDetails

	for id, loc := range *locationsWithId {
		possibleLocationsCoordinates = append(possibleLocationsCoordinates, model.PossibleLocationDetails{
			Id:          &id,
			Coordinates: loc.Coordinates,
		})
	}

	routeRequest := model.RouteToLocation{
		SourceCoordinates: &model.Coordinates{
			Longitude: *sourceLongitude,
			Latitude:  *sourceLatitude,
		},
		PossibleLocationsDetails: &possibleLocationsCoordinates,
	}

	postBody, _ := json.Marshal(routeRequest)
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(config.AppConfig.RouteServiceUri, "application/json", responseBody)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var timesResponse *[]model.LocationsTimes
	err = json.NewDecoder(resp.Body).Decode(&timesResponse)
	if err != nil {
		return nil, err
	}

	return timesResponse, nil
}

func (s *service) getAdditionalLocations(additionalLocation *model.LocationByDistanceRequestDetails, mainLocationsIds *[]int64) (*[]model.Location, *float64, error) {

	additionalLocationsQueryData := optimizeDistanceQueryParameters(additionalLocation.RadiusStart, additionalLocation.RadiusEnd)
	additionalLocationsQueryData.Type = *additionalLocation.Type

	var additionalLocationFinalRadiusEnd float64
	var additionalLocations []model.Location

	for additionalLocationsQueryData.RadiusEnd <= *additionalLocation.RadiusEnd {

		additionalLocationFinalRadiusEnd = additionalLocationsQueryData.RadiusEnd

		additionalLocationsResult, additionalLocationsErr := s.repository.GetAdditionalLocations(additionalLocationsQueryData, mainLocationsIds)

		if additionalLocationsErr != nil {
			return nil, nil, additionalLocationsErr
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

		logrus.Info("Additional locations end radius: ", additionalLocationsQueryData.RadiusEnd)
		logrus.Info("Additional locations count: ", len(additionalLocations))

		if len(additionalLocations) != 0 || *additionalLocation.RadiusEnd == additionalLocationsQueryData.RadiusEnd {
			break
		}

		if *additionalLocation.RadiusEnd-additionalLocationsQueryData.RadiusEnd > 7500 {
			additionalLocationsQueryData.RadiusEnd += 7500
		} else {
			additionalLocationsQueryData.RadiusEnd += *additionalLocation.RadiusEnd - additionalLocationsQueryData.RadiusEnd
		}
	}

	if len(additionalLocations) == 0 {
		return &[]model.Location{}, &additionalLocationFinalRadiusEnd, nil
	}

	return &additionalLocations, &additionalLocationFinalRadiusEnd, nil
}

func optimizeDistanceQueryParameters(radiusStart *float64, radiusEnd *float64) *model.LocationQuery {
	var optimizedRequest model.LocationQuery

	if *radiusEnd > 7500 && *radiusStart == 0 {
		optimizedRequest.RadiusEnd = 2000 * math.Log10(*radiusEnd)
		optimizedRequest.RadiusStart = 0
		logrus.Info("Optimizing end radius, start radius == 0")
	} else if *radiusEnd > 7500 && *radiusStart > 7500 {
		if *radiusEnd-*radiusStart < 7500 {
			optimizedRequest.RadiusEnd = *radiusEnd
			optimizedRequest.RadiusStart = *radiusStart
			logrus.Info("Original radius bigger than 7500")
		} else {
			optimizedRequest.RadiusEnd = 2000 * math.Log10(*radiusEnd)
			optimizedRequest.RadiusStart = *radiusStart
			logrus.Info("Optimizing end radius, start radius != 0")
		}
	} else {
		optimizedRequest.RadiusEnd = *radiusEnd
		optimizedRequest.RadiusStart = *radiusStart
		logrus.Info("Original radius")
	}

	return &optimizedRequest
}
