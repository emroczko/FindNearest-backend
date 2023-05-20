package util

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func APIResponse(ctx *gin.Context, StatusCode int, Data interface{}) {

	if StatusCode >= 400 {
		ctx.JSON(StatusCode, Data)
		defer ctx.AbortWithStatus(StatusCode)
	} else {
		ctx.JSON(StatusCode, Data)
	}
}

func CreateErrorResponse(ctx *gin.Context, statusCode int, error error) {

	var message string
	if statusCode == 500 {
		message = "Internal server error"
	} else {
		message = error.Error()
	}

	errResponse := ErrorResponse{
		Message: message,
	}

	ctx.JSON(statusCode, errResponse)
	defer ctx.AbortWithStatus(statusCode)
}
