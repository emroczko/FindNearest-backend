package locations

import (
	"bytes"
	"databaseClient/config"
	"databaseClient/model"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"math"
	"net/http"
)

type Service interface {
	GetLocationsByRadius(input *model.LocationByDistanceRequest) (*model.PagedLocation, error)
	GetLocationsByTime(input *model.LocationByTimeRequest) (*model.PagedLocation, error)
}

type service struct {
	repository Repository
}

func NewServiceResult(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) GetLocationsByRadius(request *model.LocationByDistanceRequest) (*model.PagedLocation, error) {

	mainLocationsQueryData := optimizeDistanceQueryParameters(request.MainLocation.RadiusStart, request.MainLocation.RadiusEnd)
	mainLocationsQueryData.Type = *request.MainLocation.Type
	mainLocationsQueryData.Longitude = *request.Longitude
	mainLocationsQueryData.Latitude = *request.Latitude
	mainLocationsQueryData.OriginalRadiusEnd = *request.MainLocation.RadiusEnd
	mainLocationsQueryData.OriginalRadiusStart = *request.MainLocation.RadiusStart

	mainLocations, mainLocationsEndRadius, mainLocationsIds, err := s.getMainLocations(mainLocationsQueryData)

	if err != nil {
		return nil, err
	}

	var pagedResults model.PagedLocation

	pagedResults.MainLocationsRadiusEnd = mainLocationsEndRadius
	mainLocationsCount := 0
	additionalLocationsRadiusEnd := 0.0
	additionalLocationsCount := 0
	if len(mainLocations) == 0 {
		pagedResults.MainLocationsCount = &mainLocationsCount
		pagedResults.MainLocations = &[]model.Location{}
		pagedResults.AdditionalLocationsCount = &additionalLocationsCount
		pagedResults.AdditionalLocationsRadiusEnd = &additionalLocationsRadiusEnd
		pagedResults.AdditionalLocations = &[]model.Location{}
	} else {
		sourceCoordinates := &model.Coordinates{
			Longitude: *request.Longitude,
			Latitude:  *request.Latitude,
		}
		locationsTimes, routeClientErr := getTimeAndDistanceToLocations(sourceCoordinates, &mainLocations, request.MeanOfTransport)

		if routeClientErr != nil {
			return nil, routeClientErr
		}

		var checkedMainLocations []model.Location
		var checkedMainLocationsIds []int64

		for _, loc := range *locationsTimes {
			checkedMainLocationsIds = append(checkedMainLocationsIds, loc.LocationDetails.Id)
			checkedMainLocations = append(checkedMainLocations, model.Location{
				Coordinates:  loc.LocationDetails.Coordinates,
				Name:         mainLocations[loc.LocationDetails.Id].Name,
				LocationType: mainLocations[loc.LocationDetails.Id].LocationType,
				Distance:     loc.Distance,
				Time:         loc.Time,
			})
		}

		pagedResults.MainLocations = &checkedMainLocations
		mainLocationsCount = len(checkedMainLocations)
		pagedResults.MainLocationsCount = &mainLocationsCount

		if request.AdditionalLocation != nil {

			if request.AdditionalLocation.Type == request.MainLocation.Type {
				return &pagedResults, nil
			}

			additionalLocations, additionalLocationsEndRadius, additionalLocErr := s.getAdditionalLocations(request.AdditionalLocation, mainLocationsIds)

			if additionalLocErr != nil {
				return nil, additionalLocErr
			}

			pagedResults.AdditionalLocations = additionalLocations
			pagedResults.AdditionalLocationsRadiusEnd = additionalLocationsEndRadius
			pagedResults.AdditionalLocationsCount = &additionalLocationsCount
		} else {
			pagedResults.AdditionalLocationsCount = &additionalLocationsCount
			pagedResults.AdditionalLocationsRadiusEnd = &additionalLocationsRadiusEnd
			pagedResults.AdditionalLocations = &[]model.Location{}
		}
	}
	return &pagedResults, nil
}

func (s *service) getMainLocations(query *model.LocationQuery) (map[int64]model.Location, *float64, *[]int64, error) {

	var pagedResults model.PagedLocation
	var mainLocationsIds []int64
	mainLocations := make(map[int64]model.Location)

	for query.RadiusEnd <= query.RadiusEnd {
		mainLocationsResult, mainLocationsErr := s.repository.GetMainLocations(query)

		if mainLocationsErr != nil {
			return nil, nil, nil, mainLocationsErr
		}

		pagedResults.MainLocationsRadiusEnd = &query.RadiusEnd

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

			var name string
			if loc.Name != nil {
				name = *loc.Name
			} else {
				name = ""
			}

			location := model.Location{
				Coordinates:  loc.Coordinates,
				Name:         &name,
				LocationType: &locationType,
			}
			mainLocations[*loc.Osm_Id] = location
			mainLocationsIds = append(mainLocationsIds, *loc.Osm_Id)
		}

		logrus.Info("End radius: ", query.RadiusEnd)
		mainLocationsCount := len(mainLocations)
		logrus.Info("Main locations length: ", mainLocationsCount)

		if len(mainLocations) != 0 || query.OriginalRadiusEnd == query.RadiusEnd {
			break
		}

		if query.OriginalRadiusEnd-query.RadiusEnd > 7500 {
			query.RadiusStart = query.RadiusEnd
			query.RadiusEnd += 7500
		} else {
			query.RadiusStart = query.RadiusEnd
			query.RadiusEnd += query.OriginalRadiusEnd - query.RadiusEnd
		}
	}

	return mainLocations, &query.RadiusEnd, &mainLocationsIds, nil
}

func getScaler(meanOfTransport *string) int {
	switch *meanOfTransport {
	case "car":
		return 1000 + 200
	case "bike":
		return 267 + 150
	case "foot":
		return 83 + 100
	default:
		return 500
	}
}

func (s *service) GetLocationsByTime(request *model.LocationByTimeRequest) (*model.PagedLocation, error) {

	scale := getScaler(request.MeanOfTransport)
	var radiusStart = float64(*request.MainLocation.TimeStart * int64(scale))
	var radiusEnd = float64(*request.MainLocation.TimeEnd * int64(scale))

	mainLocationsQueryData := optimizeDistanceQueryParameters(&radiusStart, &radiusEnd)
	mainLocationsQueryData.Type = *request.MainLocation.Type
	mainLocationsQueryData.Longitude = *request.Longitude
	mainLocationsQueryData.Latitude = *request.Latitude
	mainLocationsQueryData.OriginalRadiusEnd = radiusEnd
	mainLocationsQueryData.OriginalRadiusStart = radiusStart

	var pagedResults model.PagedLocation

	mainLocations, mainLocationsRadiusEnd, _, err := s.getMainLocations(mainLocationsQueryData)

	if err != nil {
		return nil, err
	}

	mainLocationsCount := 0
	pagedResults.MainLocationsRadiusEnd = mainLocationsRadiusEnd
	additionalLocationsRadiusEnd := 0.0
	additionalLocationsCount := 0
	if len(mainLocations) == 0 {
		pagedResults.MainLocationsCount = &mainLocationsCount
		pagedResults.MainLocations = &[]model.Location{}
		pagedResults.AdditionalLocationsCount = &additionalLocationsCount
		pagedResults.AdditionalLocationsRadiusEnd = &additionalLocationsRadiusEnd
		pagedResults.AdditionalLocations = &[]model.Location{}
	} else {
		sourceCoordinates := &model.Coordinates{
			Longitude: *request.Longitude,
			Latitude:  *request.Latitude,
		}
		locationsTimes, timesErr := getTimeAndDistanceToLocations(sourceCoordinates, &mainLocations, request.MeanOfTransport)

		if timesErr != nil {
			return nil, timesErr
		}

		var checkedMainLocations []model.Location
		var checkedMainLocationsIds []int64

		for _, loc := range *locationsTimes {
			if *loc.Time/60000 <= *request.MainLocation.TimeEnd {
				checkedMainLocationsIds = append(checkedMainLocationsIds, loc.LocationDetails.Id)
				checkedMainLocations = append(checkedMainLocations, model.Location{
					Coordinates:  loc.LocationDetails.Coordinates,
					Name:         mainLocations[loc.LocationDetails.Id].Name,
					LocationType: mainLocations[loc.LocationDetails.Id].LocationType,
					Distance:     loc.Distance,
					Time:         loc.Time,
				})
			}
		}

		pagedResults.MainLocations = &checkedMainLocations
		mainLocationsCount = len(checkedMainLocations)
		pagedResults.MainLocationsCount = &mainLocationsCount

		if request.AdditionalLocation != nil {

			if request.AdditionalLocation.Type == request.MainLocation.Type {
				return &pagedResults, nil
			}

			additionalLocations, additionalLocationsEndRadius, err := s.getAdditionalLocations(request.AdditionalLocation, &checkedMainLocationsIds)

			if err != nil {
				return nil, err
			}

			pagedResults.AdditionalLocations = additionalLocations
			additionalLocationsCount = len(*additionalLocations)
			pagedResults.AdditionalLocationsCount = &additionalLocationsCount
			pagedResults.AdditionalLocationsRadiusEnd = additionalLocationsEndRadius

		} else {
			pagedResults.AdditionalLocationsCount = &additionalLocationsCount
			pagedResults.AdditionalLocationsRadiusEnd = &additionalLocationsRadiusEnd
			pagedResults.AdditionalLocations = &[]model.Location{}
		}

	}
	return &pagedResults, nil
}

func getTimeAndDistanceToLocations(sourceCoordinates *model.Coordinates, locationsWithId *map[int64]model.Location, meanOfTransport *string) (*[]model.LocationsRouteDetails, error) {

	var possibleLocationsCoordinates []model.PossibleLocationDetails

	for id, loc := range *locationsWithId {
		possibleLocationsCoordinates = append(possibleLocationsCoordinates, model.PossibleLocationDetails{
			Id:          id,
			Coordinates: loc.Coordinates,
		})
	}

	routeRequest := model.RouteToLocation{
		SourceCoordinates:        sourceCoordinates,
		PossibleLocationsDetails: &possibleLocationsCoordinates,
		MeanOfTransport:          meanOfTransport,
	}

	postBody, _ := json.Marshal(routeRequest)
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(config.AppConfig.RouteServiceUri, "application/json", responseBody)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var timesResponse *[]model.LocationsRouteDetails
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
			additionalLocationsQueryData.RadiusStart = additionalLocationsQueryData.RadiusEnd
			additionalLocationsQueryData.RadiusEnd += 7500
		} else {
			additionalLocationsQueryData.RadiusStart = additionalLocationsQueryData.RadiusEnd
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
