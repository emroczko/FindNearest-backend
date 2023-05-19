package locationsHandler

import (
	"databaseClient/config"
	resultLocation "databaseClient/controllers/locations"
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

	//input.Distance = ctx.GetFloat64("distance")

	//errResponse, errCount := util.GoValidator(&input, config.Options)
	//
	//if errCount > 0 {
	//	util.ValidatorErrorResponse(ctx, http.StatusBadRequest, http.MethodGet, errResponse)
	//	return
	//}

	if err := ctx.BindJSON(&input); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	locations, err := h.service.ResultLocationsService(&input)

	if len(err) != 0 {
		util.CreateErrorResponse(ctx, http.StatusInternalServerError, err)
		fmt.Println(err)
	}

	if len(*locations) == 0 {
		util.APIResponse(ctx, http.StatusNotFound, locations)
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
