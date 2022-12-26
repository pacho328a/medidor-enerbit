package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	app "medidor_enerbit/app"
	utils "medidor_enerbit/utils"
)

func init() {
	// Set gin mode
	mode := utils.GetEnvVar("GIN_MODE")
	gin.SetMode(mode)
}

func main() {
	// Setup the app
	app := app.SetupApp()

	// Read ADDR and port
	addr := utils.GetEnvVar("GIN_ADDR")
	port := utils.GetEnvVar("GIN_PORT")

	// HTTP mode
	log.Info().Msgf("Starting service on http//:%s:%s", addr, port)
	if err := app.Run(fmt.Sprintf("%s:%s", addr, port)); err != nil {
		log.Fatal().Err(err).Msg("Error occurred while setting up the server")
	}

}
