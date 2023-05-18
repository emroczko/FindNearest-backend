package locations_types

import (
	"databaseClient/config"
)

type Service interface {
	ResultLocationsTypesService() []string
}

type service struct{}

func NewServiceResult() *service {
	return &service{}
}

func (s *service) ResultLocationsTypesService() []string {

	return config.LocationsTypesConfig.LocationTypes
}
