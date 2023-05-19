package util

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error interface{} `json:"message"`
}

func APIResponse(ctx *gin.Context, StatusCode int, Data interface{}) {

	if StatusCode >= 400 {
		ctx.JSON(StatusCode, Data)
		defer ctx.AbortWithStatus(StatusCode)
	} else {
		ctx.JSON(StatusCode, Data)
	}
}

func CreateErrorResponse(ctx *gin.Context, StatusCode int, Error interface{}) {
	errResponse := ErrorResponse{
		Error: Error,
	}

	ctx.JSON(StatusCode, errResponse)
	defer ctx.AbortWithStatus(StatusCode)
}
