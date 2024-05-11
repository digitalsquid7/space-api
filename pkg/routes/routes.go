package routes

import (
	"errors"
	"net/http"
	"space-api/docs"
	"space-api/pkg/exoplanetsapi"
	"space-api/pkg/exoplanetsrepo"
	"space-api/pkg/models"
	"space-api/pkg/sqlutil"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"

	// import the dialect
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

func CreateEngine() *gin.Engine {
	engine := gin.Default()
	engine.Use(ErrorMiddleware)
	docs.SwaggerInfo.BasePath = "/api/v1"

	engine.Use()

	connInfo := "host=localhost port=5432 user=test password=test dbname=space sslmode=disable"
	queryBuilder := sqlutil.NewQueryBuilder()
	exoplanetsRepository := exoplanetsrepo.New(queryBuilder, connInfo)
	exoplanetsHandler := exoplanetsapi.NewHandler(queryBuilder, exoplanetsRepository)

	v1 := engine.Group("/api/v1")
	v1.GET("/exoplanets", exoplanetsHandler.Get)
	v1.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	return engine
}

func ErrorMiddleware(c *gin.Context) {
	c.Next()

	if len(c.Errors) == 0 {
		return
	}

	var inputError *models.InvalidInputError
	if errors.As(c.Errors[0], &inputError) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: c.Errors[0].Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal Server Error"})
}
