package locations

import (
	"databaseClient/model"
)

type Service interface {
	ResultLocationsService(input *LocationRequest) (*[]model.Location, error)
}

type service struct {
	repository Repository
}

func NewServiceResult(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) ResultLocationsService(request *LocationRequest) (*[]model.Location, error) {

	pointLocationResults, err := s.repository.PointLocationsRepository(request)

	if err != nil {
		return nil, err
	}

	polygonsLocationsResults, err := s.repository.PolygonLocationsRepository(request)

	if err != nil {
		return nil, err
	}

	result := append(*pointLocationResults, *polygonsLocationsResults...)
	return &result, nil
}
