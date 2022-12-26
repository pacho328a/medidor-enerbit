package routers

import (
	controllers "medidor_enerbit/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouters(app *gin.Engine) {
	v1 := app.Group("/v1")
	{
		v1.POST("/medidor", controllers.CreateMedidor)
		v1.GET("/medidor/:id", controllers.GetMedidor)
		v1.GET("/medidores", controllers.GetMedidores)
		v1.PATCH("/medidor", controllers.UpdateMedidor)
		v1.DELETE("/medidor/:id", controllers.DeleteMedidor)
	}

}
