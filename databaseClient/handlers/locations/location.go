package locationsHandler

import (
	"databaseClient/config"
	resultLocation "databaseClient/controllers/locations"
	"databaseClient/model"
	"databaseClient/util"
	"fmt"
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

	var input resultLocation.LocationRequest

	if err := ctx.ShouldBindQuery(&input); err != nil {
		util.CreateErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	fmt.Println(input)
	locations, err := h.service.ResultLocationsService(&input)

	if err != nil {
		util.CreateErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	if len(locations.Locations) == 0 {
		util.APIResponse(ctx, http.StatusNotFound, []model.LocationEntity{})
	} else {
		util.APIResponse(ctx, http.StatusOK, locations)
	}
}

func resolveLocationType(requestLocationType string) string {
	for _, locationType := range config.LocationsTypesConfig.LocationTypes {
		if requestLocationType == locationType {
			return locationType
		}
	}
	return ""
}
