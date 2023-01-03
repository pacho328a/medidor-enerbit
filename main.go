package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"medidor_enerbit/app"
	utils "medidor_enerbit/utils"
)

func init() {
	// Set gin mode
	mode := utils.GetEnvVar("GIN_MODE")
	gin.SetMode(mode)
}

func main() {

	appx := fx.New(
		fx.Provide(
			app.SetupApp,
			utils.SetConfGin,
		),
		fx.Invoke(LifeCycleHook),
	)

	appx.Run()
}

func LifeCycleHook(lc fx.Lifecycle, app *gin.Engine, cg *utils.ConfGin) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				log.Info().Msgf("Starting service on http//:%s:%s", cg.ADDR, cg.PORT)
				go app.Run(fmt.Sprintf("%s:%s", cg.ADDR, cg.PORT))
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return nil
			},
		})
}
