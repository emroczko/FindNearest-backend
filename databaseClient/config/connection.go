package config

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func Connection() *pgxpool.Pool {

	connectionString := "postgres://" + AppConfig.DBUser + ":" + AppConfig.DBPassword + "@" + AppConfig.DBHost + ":" + AppConfig.DBPort + "/" + AppConfig.DBDatabase
	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		defer logrus.Error("Unable to parse DATABASE_URL")
		logrus.Fatal(err.Error())
	}

	var db *pgxpool.Pool
	db, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		defer logrus.Error("Unable to create connection pool")
		logrus.Fatal(err.Error())
	}

	return db
}
