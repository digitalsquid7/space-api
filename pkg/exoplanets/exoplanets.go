package exoplanets

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// Get godoc
// @Router       /exoplanets [get]
// @Summary      Read a list of exoplanets
// @Tags         exoplanets
// @Accept       json
// @Produce      json
// @Success      200  {object}  Exoplanet
func (h *Handler) Get(ctx *gin.Context) {
	response := Response[Exoplanet]{
		Data: []Exoplanet{
			{
				Id:         1,
				PlanetName: "Exo a",
				HostName:   "Exo",
			},
			{
				Id:         2,
				PlanetName: "Exo b",
				HostName:   "Exo",
			},
		},
	}

	ctx.JSON(http.StatusOK, &response)
}
