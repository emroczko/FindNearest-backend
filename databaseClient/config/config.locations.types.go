package config

import (
	"databaseClient/util/constants"
	"github.com/spf13/viper"
	"log"
)

var LocationsTypesConfig AppData

type AppData struct {
	LocationTypes []string `mapstructure:"locationTypes"`
}

func InitializeKnownLocationTypes() error {
	viper.AddConfigPath("./config/")
	viper.SetConfigName("locations.yaml")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err.Error())
		return constants.LoadConfigError
	}

	err = viper.Unmarshal(&LocationsTypesConfig)
	if err != nil {
		return constants.ParseKnownLocationsError
	}

	if len(LocationsTypesConfig.LocationTypes) == 0 {
		return constants.LoadLocationsError
	}

	return nil
}
