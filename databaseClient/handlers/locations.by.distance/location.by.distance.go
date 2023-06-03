package locationsHandler

import (
	"databaseClient/config"
	resultLocation "databaseClient/controllers/locations"
	"databaseClient/model"
	"databaseClient/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	service resultLocation.Service
}

func NewHandlerLocationByDistance(service resultLocation.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetLocationsByDistance(ctx *gin.Context) {

	var input *model.LocationByDistanceRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		util.CreateErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	logrus.Info("Request: ", *input.Latitude, *input.Longitude, *input.MainLocation.Type)
	locations, err := h.service.GetLocationsByDistance(input)

	if err != nil {
		util.CreateErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	if len(*locations.MainLocations) == 0 {
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
