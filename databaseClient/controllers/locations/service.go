package locations

import (
	"databaseClient/model"
)

type Service interface {
	ResultLocationsService(input *model.LocationRequest) (*model.PagedLocation, error)
}

type service struct {
	repository Repository
}

func NewServiceResult(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) ResultLocationsService(request *model.LocationRequest) (*model.PagedLocation, error) {

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

	pagedResults.RadiusEnd = &request.RadiusEnd
	var locations []model.Location

	for _, loc := range result {
		location := model.Location{
			Coordinates:  loc.Coordinates,
			Name:         loc.Name,
			LocationType: loc.Amenity,
		}
		locations = append(locations, location)
	}

	pagedResults.Locations = &locations

	return &pagedResults, nil
}
