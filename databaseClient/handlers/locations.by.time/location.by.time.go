package locations_by_time

import (
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

func NewHandlerLocationByTime(service resultLocation.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetLocationsByTime(ctx *gin.Context) {

	var input *model.LocationByTimeRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		util.CreateErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	logrus.Info("Request: ", *input.Latitude, *input.Longitude, *input.MainLocation.Type)
	locations, err := h.service.GetLocationsByTime(input)

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
