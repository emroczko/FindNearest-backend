package config

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func Connection() *pgxpool.Pool {

	connectionString := "postgres://" + AppConfig.DBUser + ":" + AppConfig.DBPassword + "@" + AppConfig.DBHost + ":" + AppConfig.DBPort + "/" + AppConfig.DBDatabase
	poolConfig, err := pgxpool.ParseConfig(connectionString)
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
