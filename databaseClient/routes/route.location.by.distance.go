package routes

import (
	location "databaseClient/controllers/locations"
	locationsHandler "databaseClient/handlers/locations"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitLocationRoutes(route *gin.Engine, conn *pgxpool.Pool) {

	resultsLocationsRepository := location.CreateRepository(conn)
	resultsLocationsByDistanceService := location.NewServiceResult(resultsLocationsRepository)
	resultsLocationsByDistanceHandler := locationsHandler.NewHandlerResultLocation(resultsLocationsByDistanceService)

	resultsLocationsByTimeService := location.NewServiceResult(resultsLocationsRepository)
	resultsLocationsByTimeHandler := locationsHandler.NewHandlerResultLocation(resultsLocationsByDistanceService)

	groupRoute := route.Group("/api/v1")
	groupRoute.POST("/locationsByDistance", resultsLocationsByDistanceHandler.ResultLocationHandler)
	groupRoute.POST("/locationsByTime", resultsLocationsByDistanceHandler.ResultLocationHandler)
}
