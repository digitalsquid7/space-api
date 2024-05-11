package exoplanetsapi

import (
	"net/http"
	"space-api/pkg/exoplanetsrepo"
	"space-api/pkg/models"
	"space-api/pkg/sqlutil"

	"github.com/gin-gonic/gin"
)

type Selector[T any] interface {
	Select() ([]T, error)
}

type Handler struct {
	fields       sqlutil.Fields
	queryBuilder *sqlutil.QueryBuilder
	repository   *exoplanetsrepo.ExoplanetsRepository
}

func NewHandler(queryBuilder *sqlutil.QueryBuilder, repository *exoplanetsrepo.ExoplanetsRepository) *Handler {
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
	queryModifiers, err := sqlutil.LoadQueryModifiers(ctx.Request,
		sqlutil.WithPaging(20),
		sqlutil.WithFilters(h.fields),
		sqlutil.WithSorting(h.fields))
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := models.Response[models.Exoplanet]{}
	if response.Data, err = h.repository.ReadExoplanets(queryModifiers); err != nil {
		_ = ctx.Error(err)
		return
	}

	if queryModifiers.Page.IncludeTotalSize {
		response.Meta = &models.Meta{}
		if response.Meta.TotalSize, err = h.repository.ReadCount(queryModifiers); err != nil {
			_ = ctx.Error(err)
			return
		}
	}

	ctx.JSON(http.StatusOK, &response)
}

func getFields() sqlutil.Fields {
	return sqlutil.Fields{
		{
			SQLName: "id",
			APIName: "id",
			Type:    sqlutil.Integer,
		},
		{
			SQLName: "planet_name",
			APIName: "planetName",
			Type:    sqlutil.String,
		},
		{
			SQLName: "host_name",
			APIName: "hostName",
			Type:    sqlutil.String,
		},
		{
			SQLName: "system_number",
			APIName: "systemNumber",
			Type:    sqlutil.Integer,
		},
		{
			SQLName: "discovery_method",
			APIName: "discoveryMethod",
			Type:    sqlutil.String,
		},
		{
			SQLName: "year_discovered",
			APIName: "yearDiscovered",
			Type:    sqlutil.String,
		},
	}
}
