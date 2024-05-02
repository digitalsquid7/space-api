package exoplanetsapi

import (
	"net/http"
	"space-api/pkg/exoplanetsrepo"
	"space-api/pkg/models"
	"space-api/pkg/querybuilder"
	"space-api/pkg/querymodifiers"

	"github.com/gin-gonic/gin"
)

type Selector[T any] interface {
	Select() ([]T, error)
}

type Handler struct {
	fields       querymodifiers.Fields
	queryBuilder *querybuilder.QueryBuilder
	repository   *exoplanetsrepo.ExoplanetsRepository
}

func NewHandler(queryBuilder *querybuilder.QueryBuilder, repository *exoplanetsrepo.ExoplanetsRepository) *Handler {
	return &Handler{
		fields:       getFields(),
		queryBuilder: queryBuilder,
		repository:   repository,
	}
}

// Get godoc
// @Router       /exoplanetsapi [get]
// @Summary      Read a list of exoplanetsapi
// @Tags         exoplanetsapi
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Exoplanet
func (h *Handler) Get(ctx *gin.Context) {
	queryModifiers, err := querymodifiers.Load(ctx.Request,
		querymodifiers.WithPaging(20),
		querymodifiers.WithFilters(h.fields),
		querymodifiers.WithSorting(h.fields))
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	exoplanets, err := h.repository.Read(queryModifiers)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := models.Response[models.Exoplanet]{Data: exoplanets}
	ctx.JSON(http.StatusOK, &response)
}

func getFields() querymodifiers.Fields {
	return querymodifiers.Fields{
		{
			SQLName: "id",
			APIName: "id",
			Type:    querymodifiers.Integer,
		},
		{
			SQLName: "planet_name",
			APIName: "planetName",
			Type:    querymodifiers.String,
		},
		{
			SQLName: "host_name",
			APIName: "hostName",
			Type:    querymodifiers.String,
		},
		{
			SQLName: "system_number",
			APIName: "systemNumber",
			Type:    querymodifiers.Integer,
		},
		{
			SQLName: "discovery_method",
			APIName: "discoveryMethod",
			Type:    querymodifiers.String,
		},
		{
			SQLName: "year_discovered",
			APIName: "yearDiscovered",
			Type:    querymodifiers.String,
		},
	}
}
