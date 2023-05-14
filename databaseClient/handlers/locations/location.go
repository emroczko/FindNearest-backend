package locationsHandler

import (
	resultLocation "databaseClient/controllers/locations"
	"databaseClient/util"
	"github.com/gin-gonic/gin"
	"log"
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

	resultStudent, err := h.service.ResultLocationsService(&input)

	if len(err) != 0 {
		log.Fatal(err)
	}

	//switch errResultStudent {
	//
	//case "RESULT_STUDENT_NOT_FOUND_404":
	//	util.APIResponse(ctx, "Student data is not exist or deleted", http.StatusNotFound, http.MethodGet, nil)
	//	return
	//
	//default:
	util.APIResponse(ctx, "Result Student data successfully", http.StatusOK, http.MethodGet, resultStudent)
	//}
}
