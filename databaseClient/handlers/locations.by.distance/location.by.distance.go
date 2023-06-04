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
		logrus.Error(&ctx.Request.Body)
		logrus.Error(err)
		util.CreateErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	logrus.Info("Request: ", *input.Latitude, *input.Longitude, *input.MainLocation.Type)
	locations, err := h.service.GetLocationsByDistance(input)

	if err != nil {
		logrus.Error(err)
		util.CreateErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	util.APIResponse(ctx, http.StatusOK, locations)
}

func resolveLocationType(requestLocationType string) string {
	for _, locationType := range config.LocationsTypesConfig.LocationTypes {
		if requestLocationType == locationType {
			return locationType
		}
	}
	return ""
}
