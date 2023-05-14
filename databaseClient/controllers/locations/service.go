package locations

import "databaseClient/model"

type Service interface {
	ResultLocationsService(input *LocationRequest) (*[]model.Location, string)
}

type service struct {
	repository Repository
}

func NewServiceResult(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) ResultLocationsService(input *LocationRequest) (*[]model.Location, string) {

	locationsResults, err := s.repository.ResultLocationsRepository(input)

	return locationsResults, err
}
