package locations

import (
	"context"
	"databaseClient/model"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	PointLocationsRepository(input *LocationRequest) (*[]model.Location, string)
	PolygonLocationsRepository(input *LocationRequest) (*[]model.Location, string)
}

type repository struct {
	conn *pgxpool.Pool
}

func NewRepositoryResult(conn *pgxpool.Pool) *repository {
	return &repository{conn: conn}
}

func (r *repository) PointLocationsRepository(input *LocationRequest) (*[]model.Location, string) {

	var locationsResult []model.Location

	sql := createQuery("PLANET_OSM_POINT")

	rows, _ := r.conn.Query(context.Background(), sql, input.Type, input.Longitude, input.Latitude, input.RadiusStart, input.RadiusEnd)

	for rows.Next() {
		location, err := pgx.RowToAddrOfStructByName[model.Location](rows)
		if err != nil {
			return &locationsResult, err.Error()
		}

		locationsResult = append(locationsResult, *location)
	}

	return &locationsResult, ""
}

func (r *repository) PolygonLocationsRepository(input *LocationRequest) (*[]model.Location, string) {

	var locationsResult []model.Location

	sql := createQuery("PLANET_OSM_POLYGON")

	rows, _ := r.conn.Query(context.Background(), sql, input.Type, input.Longitude, input.Latitude, input.RadiusStart, input.RadiusEnd)

	for rows.Next() {
		location, err := pgx.RowToAddrOfStructByName[model.Location](rows)
		if err != nil {
			return &locationsResult, err.Error()
		}

		locationsResult = append(locationsResult, *location)
	}

	return &locationsResult, ""
}

func createQuery(tableName string) string {
	sql := fmt.Sprintf(`
		SELECT AMENITY,
			NAME,
			SHOP,
			SPORT,
			PUBLIC_TRANSPORT,
			TAGS,
			WATER,
			LANDUSE
		FROM %s
		WHERE HIGHWAY IS NULL
			AND RAILWAY IS NULL
			AND POWER IS NULL
			AND BARRIER IS NULL
			AND (BUILDING != 'garage'
								AND BUILDING != 'apartments'
								OR BUILDING IS NULL)
			AND (LANDUSE != 'grass'
								OR LANDUSE IS NULL)
			AND (SHOP = $1
								OR LEISURE = $1
								OR PUBLIC_TRANSPORT = $1
			         			OR AMENITY = $1
								OR WATER = $1)
			AND NOT ((NAME IS NULL OR NAME = '') AND (TAGS IS NULL OR TAGS = ''))
			AND ($4 = 0 OR ST_DWITHIN(WAY, ST_TRANSFORM(ST_SETSRID(ST_POINT($2, $3), 4326), 3857), $4) = false)
			AND ST_DWITHIN(WAY, ST_TRANSFORM(ST_SETSRID(ST_POINT($2, $3), 4326), 3857), $5)
	`, tableName)

	return sql
}
