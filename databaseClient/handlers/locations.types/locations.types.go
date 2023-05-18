package locationsTypesHandler

import (
	resultLocation "databaseClient/controllers/locations.types"
	"databaseClient/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	service resultLocation.Service
}

func NewHandlerResultLocation(service resultLocation.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ResultLocationHandler(ctx *gin.Context) {
	locationsTypes := h.service.ResultLocationsTypesService()
	util.APIResponse(ctx, "Locations types data returned successfully", http.StatusOK, http.MethodGet, locationsTypes)
}
