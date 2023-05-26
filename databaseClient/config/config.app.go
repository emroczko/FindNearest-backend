package config

import (
	"databaseClient/util/constants"
	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	Port        string `mapstructure:"GO_PORT"`
	Host        string `mapstructure:"GO_HOST"`
	Environment string `mapstructure:"GO_ENV"`

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

//func GetFromEnv() error {
//
//	AppConfig.Host = util.GetEnv("GO_HOST")
//	AppConfig.Port = util.GetEnv("GO_PORT")
//	AppConfig.Environment = util.GetEnv("GO_PORT")
//
//	AppConfig.DBHost = util.GetEnv("DB_PSQL_HOST")
//	AppConfig.DBPort = util.GetEnv("DB_PSQL_PORT")
//	AppConfig.DBUser = util.GetEnv("DB_PSQL_USERNAME")
//	AppConfig.DBPassword = util.GetEnv("DB_PSQL_PASSWORD")
//	AppConfig.DBDatabase = util.GetEnv("DB_PSQL_DATABASE")
//
//	validate = validator.New()
//
//	err := validate.Struct(AppConfig)
//	if err != nil {
//		return err
//	}
//}
