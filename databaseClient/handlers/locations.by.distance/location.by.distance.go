package locationsDistanceHandler

import (
	resultLocation "databaseClient/controllers/locations"
	"databaseClient/model"
	"databaseClient/util"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
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

	if resolveMeanOfTransport(input.MeanOfTransport) == false {
		util.CreateErrorResponse(ctx, http.StatusBadRequest, errors.New("not allowed mean of transport"))
		return
	}

	logrus.Info("Request: ",
		" latitude: ", *input.Latitude,
		" longitude: ", *input.Longitude,
		" type: ", *input.MainLocation.Type,
		" mean of transport: ", *input.MeanOfTransport)
	
	locations, err := h.service.GetLocationsByRadius(input)

	if err != nil {
		logrus.Error(err)
		util.CreateErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	util.APIResponse(ctx, http.StatusOK, locations)
}

func resolveMeanOfTransport(requestLocationType *string) bool {
	allowedMeansOfTransport := []string{"car", "foot", "bike"}
	return slices.Contains(allowedMeansOfTransport, *requestLocationType)
}
