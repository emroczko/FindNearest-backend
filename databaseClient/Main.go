package main

import (
	"databaseClient/config"
	"databaseClient/routes"
	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := setupRouter()
	log.Fatal(router.Run("localhost:8080"))
}

func setupRouter() *gin.Engine {

	db := config.Connection()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))
	router.Use(helmet.Default())
	router.Use(gzip.Gzip(gzip.BestCompression))

	routes.InitLocationRoutes(router, db)

	return router
}
