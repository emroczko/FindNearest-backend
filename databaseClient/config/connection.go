package config

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func Connection() *pgxpool.Pool {

	poolConfig, err := pgxpool.ParseConfig("")
	if err != nil {
		log.Fatalln("Unable to parse DATABASE_URL:", err)
	}

	var db *pgxpool.Pool
	db, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalln("Unable to create connection pool:", err)
	}

	return db
}
