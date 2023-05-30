package locations

import (
	"context"
	"databaseClient/model"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	PointLocationsRepository(input *model.LocationRequest) (*[]model.LocationEntity, error)
	PolygonLocationsRepository(input *model.LocationRequest) (*[]model.LocationEntity, error)
}

type repository struct {
	conn *pgxpool.Pool
}

func NewRepositoryResult(conn *pgxpool.Pool) *repository {
	return &repository{conn: conn}
}

func (r *repository) PointLocationsRepository(input *model.LocationRequest) (*[]model.LocationEntity, error) {

	var locationsResult []model.LocationEntity

	sql := createQuery(POINTS)

	rows, err := r.conn.Query(context.Background(), sql, input.Type, input.Longitude, input.Latitude, input.RadiusStart, input.RadiusEnd)

	if err != nil {
		logrus.Error("Database error:", err.Error())
		return nil, err
	}

	for rows.Next() {
		location, err := pgx.RowToAddrOfStructByName[model.LocationEntity](rows)
		if err != nil {
			logrus.Error("Parsing database data error:", err.Error())
			return &locationsResult, err
		}

		locationsResult = append(locationsResult, *location)
	}

	return &locationsResult, nil
}

func (r *repository) PolygonLocationsRepository(input *model.LocationRequest) (*[]model.LocationEntity, error) {

	var locationsResult []model.LocationEntity

	sql := createQuery(POLYGONS)

	rows, err := r.conn.Query(context.Background(), sql, input.Type, input.Longitude, input.Latitude, input.RadiusStart, input.RadiusEnd)

	if err != nil {
		logrus.Error("Database error: ", err.Error())
		return nil, err
	}

	for rows.Next() {
		location, err := pgx.RowToAddrOfStructByName[model.LocationEntity](rows)
		if err != nil {
			logrus.Error("Parsing database data error: ", err.Error())
			return &locationsResult, err
		}

		locationsResult = append(locationsResult, *location)
	}

	return &locationsResult, err
}

func createQuery(tableName TABLE) string {

	var geometryColumn string

	if tableName == POLYGONS {
		geometryColumn = "ST_AsText(ST_PointN(ST_Exteriorring(ST_Transform(WAY, 4326)), 1)) AS COORDINATES"
	} else {
		geometryColumn = "ST_AsText(ST_Transform(WAY, 4326)) AS COORDINATES"
	}

	sql := fmt.Sprintf(`
		SELECT AMENITY,
			NAME,
			SHOP,
			SPORT,
			PUBLIC_TRANSPORT,
			TAGS,
			WATER,
			LANDUSE,
			%s
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
	`, geometryColumn, tableName)

	return sql
}

type TABLE string

const (
	POINTS   TABLE = "PLANET_OSM_POINT"
	POLYGONS TABLE = "PLANET_OSM_POLYGON"
)
