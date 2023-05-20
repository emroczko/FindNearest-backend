package config

import (
	constants "databaseClient/util/constants"
	"github.com/spf13/viper"
	"log"
)

var AppConfig Config

type Config struct {
	Port        int    `mapstructure:"PORT"`
	Host        string `mapstructure:"HOST"`
	Environment string `mapstructure:"ENVIRONMENT"`

	DBHost     string `mapstructure:"DB_PSQL_HOST"`
	DBPort     string `mapstructure:"DB_PSQL_PORT"`
	DBUser     string `mapstructure:"DB_PSQL_USERNAME"`
	DBPassword string `mapstructure:"DB_PSQL_PASSWORD"`
	DBDatabase string `mapstructure:"DB_PSQL_DATABASE"`
}

func InitializeAppConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err.Error())
		return constants.LoadEnvsError
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		return constants.ParseConfigError
	}

	if AppConfig.Port == 0 || AppConfig.Environment == "" {
		return constants.EmptyEnvVarError
	}

	if AppConfig.DBUser == "" || AppConfig.DBPassword == "" {
		return constants.EmptyEnvVarError
	}

	return nil
}
