package utils

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type ConfGin struct {
	ADDR string
	PORT string
}

func init() {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Debug().Err(err).
			Msg("Error occurred while reading env file, might fallback to OS env config")
	}
	viper.AutomaticEnv()
}

// This function can be used to get ENV Var in our App
// Modify this if you change the library to read ENV
func GetEnvVar(name string) string {
	if !viper.IsSet(name) {
		log.Debug().Msgf("Environment variable %s is not set", name)
		return ""
	}
	value := viper.GetString(name)
	return value
}

func SetConfGin() *ConfGin {
	addr := GetEnvVar("GIN_ADDR")
	port := GetEnvVar("GIN_PORT")

	gc := ConfGin{addr, port}

	return &gc
}
