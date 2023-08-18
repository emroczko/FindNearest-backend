package config

import (
	"databaseClient/util/constants"
	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	Port        string `mapstructure:"PORT"`
	Environment string `mapstructure:"ENV"`

	DBHost     string `mapstructure:"DB_PSQL_HOST"`
	DBPort     string `mapstructure:"DB_PSQL_PORT"`
	DBUser     string `mapstructure:"DB_PSQL_USERNAME"`
	DBPassword string `mapstructure:"DB_PSQL_PASSWORD"`
	DBDatabase string `mapstructure:"DB_PSQL_DATABASE"`

	RouteServiceUri string `mapstructure:"ROUTE_SERVICE_URI"`
}

func InitializeAppConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if bindErr := viper.BindStruct(&AppConfig); bindErr != nil {
				return constants.LoadEnvsError
			}
		} else {
			return constants.ParseConfigError
		}
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		return constants.ParseConfigError
	}

	return nil
}
