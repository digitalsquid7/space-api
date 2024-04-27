package routes

import (
	"space-api/docs"
	"space-api/pkg/exoplanets"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

func CreateEngine() *gin.Engine {
	engine := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	exoplanetsHandler := exoplanets.NewHandler()

	v1 := engine.Group("/api/v1")
	v1.GET("/exoplanets", exoplanetsHandler.Get)
	v1.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	return engine
}
