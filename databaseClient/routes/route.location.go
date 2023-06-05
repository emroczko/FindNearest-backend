package routes

import (
	location "databaseClient/controllers/locations"
	locationsDistanceHandler "databaseClient/handlers/locations.by.distance"
	locationsTimeHandler "databaseClient/handlers/locations.by.time"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitLocationRoutes(route *gin.Engine, conn *pgxpool.Pool) {

	resultsLocationsRepository := location.CreateRepository(conn)
	resultsLocationsService := location.NewServiceResult(resultsLocationsRepository)
	resultsLocationsByDistanceHandler := locationsDistanceHandler.NewHandlerLocationByDistance(resultsLocationsService)
	resultsLocationsByTimeHandler := locationsTimeHandler.NewHandlerLocationByTime(resultsLocationsService)

	groupRoute := route.Group("/api/v1")
	groupRoute.POST("/locationsByDistance", resultsLocationsByDistanceHandler.GetLocationsByDistance)
	groupRoute.POST("/locationsByTime", resultsLocationsByTimeHandler.GetLocationsByTime)
}
