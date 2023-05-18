package locations

import (
	"context"
	"databaseClient/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	ResultLocationsRepository(input *LocationRequest) (*[]model.Location, string)
}

type repository struct {
	conn *pgxpool.Pool
}

func NewRepositoryResult(conn *pgxpool.Pool) *repository {
	return &repository{conn: conn}
}

func (r *repository) ResultLocationsRepository(input *LocationRequest) (*[]model.Location, string) {

	var locationsResult []model.Location

	var sql = `
		SELECT AMENITY,
			NAME,
			SHOP,
			SPORT,
			PUBLIC_TRANSPORT,
			TAGS,
			WATER,
			LANDUSE
		FROM PLANET_OSM_POINT
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
		-- 	AND ST_DWITHIN(WAY, ST_TRANSFORM(ST_SETSRID(ST_POINT(20.97954, 52.25052), 4326), 3857), 1000) = false
			AND ST_DWITHIN(WAY, ST_TRANSFORM(ST_SETSRID(ST_POINT($2, $3), 4326), 3857), $4)
	`

	rows, _ := r.conn.Query(context.Background(), sql, input.Type, input.Longitude, input.Latitude, input.Distance)

	for rows.Next() {
		location, err := pgx.RowToAddrOfStructByName[model.Location](rows)
		if err != nil {
			return &locationsResult, err.Error()
		}

		locationsResult = append(locationsResult, *location)
	}

	return &locationsResult, ""
}
