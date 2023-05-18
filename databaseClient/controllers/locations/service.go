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

	locationsResults, err := s.repository.ResultLocationsRepository(request)

	fmt.Println(request)

	return locationsResults, err
}

func resolveLocationType(requestLocationType string) string {
	for _, locationType := range config.LocationsTypesConfig.LocationTypes {
		if requestLocationType == locationType {
			return locationType
		}
	}
	return ""
}
