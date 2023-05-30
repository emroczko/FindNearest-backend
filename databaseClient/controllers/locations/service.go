package locations

import (
	"databaseClient/model"
)

type Service interface {
	ResultLocationsService(input *LocationRequest) (*model.PagedLocation, error)
}

type service struct {
	repository Repository
}

func NewServiceResult(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) ResultLocationsService(request *LocationRequest) (*model.PagedLocation, error) {

	pointLocationResults, err := s.repository.PointLocationsRepository(request)

	if err != nil {
		return nil, err
	}

	polygonsLocationsResults, err := s.repository.PolygonLocationsRepository(request)

	if err != nil {
		return nil, err
	}

	result := append(*pointLocationResults, *polygonsLocationsResults...)

	var pagedResults model.PagedLocation

	pagedResults.PageNumber = 0
	pagedResults.PagesCount = 0

	for _, loc := range result {
		location := model.Location{
			Coordinates: *loc.Coordinates,
			Name:        *loc.Name,
			Amenity:     *loc.Amenity,
		}
		pagedResults.Locations = append(pagedResults.Locations, location)
	}

	return &pagedResults, nil
}
