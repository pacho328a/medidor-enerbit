package main

import (
	"context"
	"fmt"
	"medidor_enerbit/app"
	utils "medidor_enerbit/utils"

	_ "medidor_enerbit/docs"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func init() {
	// Set gin mode
	mode := utils.GetEnvVar("GIN_MODE")
	gin.SetMode(mode)
}

// @title           Medidor Service
// @version         1.0
// @description     A medidor management service API in Go using GORM.

// @contact.name   Francisco Anacona
// @contact.url    http://artemisa.unicauca.edu.co/~javieranacona/index.html
// @contact.email  pacho328@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:5000
// @BasePath  /v1
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
