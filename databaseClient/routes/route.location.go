package routes

import (
	location "databaseClient/controllers/locations"
	locationsHandler "databaseClient/handlers/locations"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitLocationRoutes(route *gin.Engine, conn *pgxpool.Pool) {

	resultsLocationsRepository := location.NewRepositoryResult(conn)
	resultsLocationsService := location.NewServiceResult(resultsLocationsRepository)
	resultsLocationsHandler := locationsHandler.NewHandlerResultLocation(resultsLocationsService)

	groupRoute := route.Group("/api/v1")
	groupRoute.GET("/locations", resultsLocationsHandler.ResultLocationHandler)
}
