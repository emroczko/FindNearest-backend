package routes

import (
	locations_types "databaseClient/controllers/locations.types"
	locationsTypesHandler "databaseClient/handlers/locations.types"
	"github.com/gin-gonic/gin"
)

func InitLocationsTypesRoutes(route *gin.Engine) {

	resultsLocationsTypesService := locations_types.NewServiceResult()
	resultsLocationsHandler := locationsTypesHandler.NewHandlerResultLocation(resultsLocationsTypesService)

	groupRoute := route.Group("/api/v1")
	groupRoute.GET("/locationsTypes", resultsLocationsHandler.ResultLocationHandler)
}
