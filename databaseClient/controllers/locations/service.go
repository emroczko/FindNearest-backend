package locations

import (
	"databaseClient/model"
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

	result := append(*pointLocationResults, *polygonsLocationsResults...)

	return &result, err
}
