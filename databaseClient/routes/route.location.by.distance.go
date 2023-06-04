package routes

import (
	location "databaseClient/controllers/locations"
	locationsHandler "databaseClient/handlers/locations.by.distance"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitLocationRoutes(route *gin.Engine, conn *pgxpool.Pool) {

	resultsLocationsRepository := location.CreateRepository(conn)
	resultsLocationsByDistanceService := location.NewServiceResult(resultsLocationsRepository)
	resultsLocationsByDistanceHandler := locationsHandler.NewHandlerLocationByDistance(resultsLocationsByDistanceService)
	resultsLocationsByTimeHandler := locationsHandler.NewHandlerLocationByDistance(resultsLocationsByDistanceService)

	groupRoute := route.Group("/api/v1")
	groupRoute.POST("/locationsByDistance", resultsLocationsByDistanceHandler.GetLocationsByDistance)
	groupRoute.POST("/locationsByTime", resultsLocationsByTimeHandler.GetLocationsByDistance)
}
