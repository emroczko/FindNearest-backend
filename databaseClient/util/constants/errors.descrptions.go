package constants

import "errors"

var (
	LoadEnvsError            = errors.New("could not load environment variables")
	LoadConfigError          = errors.New("could not load config")
	LoadLocationsError       = errors.New("loaded locations are empty")
	ParseKnownLocationsError = errors.New("known locations parsing has failed")
	ParseConfigError         = errors.New("config parsing has failed")
	EmptyEnvVarError         = errors.New("required environment variable is empty")
)
