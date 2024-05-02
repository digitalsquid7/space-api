package routes

import (
	"space-api/docs"
	"space-api/pkg/exoplanetsapi"
	"space-api/pkg/exoplanetsrepo"
	"space-api/pkg/querybuilder"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"

	// import the dialect
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

func CreateEngine() *gin.Engine {
	engine := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	queryBuilder := querybuilder.New()
	exoplanetsRepository := exoplanetsrepo.New(queryBuilder)
	exoplanetsHandler := exoplanetsapi.NewHandler(queryBuilder, exoplanetsRepository)

	v1 := engine.Group("/api/v1")
	v1.GET("/exoplanets", exoplanetsHandler.Get)
	v1.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	return engine
}
