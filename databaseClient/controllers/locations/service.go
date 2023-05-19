package locations

import (
	"databaseClient/config"
	"databaseClient/model"
	"fmt"
)

type Service interface {
	ResultLocationsService(input *LocationRequest) (*[]model.Location, string)
}

type service struct {
	repository Repository
}

func NewServiceResult(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) ResultLocationsService(request *LocationRequest) (*[]model.Location, string) {

	pointLocationResults, err := s.repository.PointLocationsRepository(request)
	polygonsLocationsResults, err := s.repository.PolygonLocationsRepository(request)

	fmt.Println(request)
	fmt.Println(pointLocationResults)
	fmt.Println(polygonsLocationsResults)

	result := append(*pointLocationResults, *polygonsLocationsResults...)

	return &result, err
}

func resolveLocationType(requestLocationType string) string {
	for _, locationType := range config.LocationsTypesConfig.LocationTypes {
		if requestLocationType == locationType {
			return locationType
		}
	}
	return ""
}
