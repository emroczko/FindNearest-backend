package locations

import (
	"context"
	"databaseClient/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

type Repository interface {
	GetMainLocations(input *model.LocationQuery) (*[]model.LocationEntity, error)
	GetAdditionalLocations(input *model.LocationQuery, mainLocationsIds *[]int64) (*[]model.LocationEntity, error)
}

type repository struct {
	conn *pgxpool.Pool
}

func CreateRepository(conn *pgxpool.Pool) *repository {
	return &repository{conn: conn}
}

func (r *repository) GetMainLocations(input *model.LocationQuery) (*[]model.LocationEntity, error) {

	var locationsResult []model.LocationEntity

	sql := `
		SELECT OSM_ID,
		       AMENITY,
				NAME,
				SHOP,
				SPORT,
				PUBLIC_TRANSPORT,
				TAGS,
				WATER,
				LANDUSE,
				ST_AsText(ST_Transform(WAY, 4326)) AS COORDINATES
		FROM PLACES
		WHERE 
			 (SHOP = $1
							OR LEISURE = $1
							OR PUBLIC_TRANSPORT = $1
							OR AMENITY = $1
							OR WATER = $1)
			AND ($4 = 0 OR ST_DWITHIN(ST_Transform(WAY, 2180), ST_TRANSFORM(ST_SETSRID(ST_POINT($2, $3), 4326), 2180), $4) = false)
			AND ST_DWITHIN(ST_Transform(WAY, 2180), ST_TRANSFORM(ST_SETSRID(ST_POINT($2, $3), 4326), 2180), $5)
	`

	start := time.Now()
	rows, err := r.conn.Query(context.Background(), sql, input.Type, input.Longitude, input.Latitude, input.RadiusStart, input.RadiusEnd)
	elapsed := time.Since(start)
	log.Printf("Main locations.by.distance query took %s", elapsed)

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

func (r *repository) GetAdditionalLocations(input *model.LocationQuery, mainLocationsIds *[]int64) (*[]model.LocationEntity, error) {

	var locationsResult []model.LocationEntity

	sql := `
		SELECT P1.OSM_ID,
			P1.AMENITY,
			P1.NAME,
			P1.SHOP,
			P1.SPORT,
			P1.PUBLIC_TRANSPORT,
			P1.TAGS,
			P1.WATER,
			P1.LANDUSE,
			ST_AsText(ST_Transform(P1.WAY, 4326)) AS COORDINATES
		FROM PLACES P1
		INNER JOIN PLACES P2
		ON ($1 = 0 OR ST_DWITHIN(ST_Transform(P1.WAY, 2180), ST_Transform(P2.WAY, 2180), $1) = false) 
		       AND ST_DWITHIN(ST_Transform(P1.WAY, 2180), ST_Transform(P2.WAY, 2180), $2)
		WHERE
    			(P1.SHOP = $3
						OR P1.LEISURE = $3
						OR P1.PUBLIC_TRANSPORT = $3
						OR P1.AMENITY = $3
						OR P1.WATER = $3)
		AND P2.OSM_ID = ANY ($4)`

	start := time.Now()
	rows, err := r.conn.Query(context.Background(), sql, input.RadiusStart, input.RadiusEnd, input.Type, *mainLocationsIds)
	elapsed := time.Since(start)
	log.Printf("Additional places query took %s", elapsed)

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
