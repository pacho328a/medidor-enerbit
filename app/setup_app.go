package app

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	grpcs "medidor_enerbit/gRPC"
	middlewares "medidor_enerbit/middlewares"
	routers "medidor_enerbit/routers"
	"medidor_enerbit/utils"
)

// Function to setup the app object
func SetupApp() *gin.Engine {
	log.Info().Msg("Initializing service")

	// Create barebone engine
	app := gin.New()
	// Add default recovery middleware
	app.Use(gin.Recovery())

	// disabling the trusted proxy feature
	app.SetTrustedProxies(nil)

	// Add cors, request ID and request logging middleware
	log.Info().Msg("Adding cors, request id and request logging middleware")
	app.Use(middlewares.CORSMiddleware(), middlewares.RequestLogger())

	// Setup routers
	log.Info().Msg("Setting up routers")
	routers.SetupRouters(app)

	log.Info().Msg("Creating the database connection")
	// Create the database connection
	if dberr := utils.CreateDBConnection(); dberr != nil {
		log.Err(dberr).Msg("Error occurred while creating the database connection")
	}

	// Auto migrate database
	err := utils.AutoMigrateDB()
	if err != nil {
		log.Err(err).Msg("Error occurred while auto migrating database")
	}
	go grpcs.GRPCListen()
	return app
}
