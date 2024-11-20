package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"pepper-analytics-ai/internal/handlers"
	"pepper-analytics-ai/internal/services"
)

type RouterConfig struct {
	DB *sqlx.DB
}

func SetupRouter(config RouterConfig) (*gin.Engine, error) {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	// Initialize services
	plantService := services.NewPlantService(config.DB)
	fileService := services.NewFileService("/uploads")

	plantHandler := handlers.NewPlantHandler(plantService, fileService)

	// Static files
	router.LoadHTMLGlob("templates/**/*")
	router.Static("/css", "./static/css")
	router.Static("/js", "./static/js")
	router.Static("/img", "./static/img")
	router.Static("/uploads", "./uploads")

	// Favicon
	router.StaticFile("/favicon.svg", "./static/img/favicon.svg")

	// Public routes
	router.GET("/", plantHandler.HandlePlantList)

	// Plant routes
	router.GET("/plants/new", plantHandler.HandleNewPlantForm)
	router.POST("/plants/create", plantHandler.HandleCreatePlant)
	router.GET("/plants/:id/edit", plantHandler.HandleEditPlantForm)
	router.PUT("/plants/:id", plantHandler.HandleUpdatePlant)
	router.DELETE("/plants/:id", plantHandler.HandleDeletePlant)

	// routes.go
	router.GET("/plants/:id/journal", plantHandler.HandleJournal)
	router.POST("/plants/:id/journal", plantHandler.HandleCreateJournalEntry)
	router.DELETE("/plants/:id/journal/:entryId", plantHandler.HandleDeleteJournalEntry)

	router.GET("/plants/:id/journal/:entryId/edit", plantHandler.HandleEditJournalEntry)
	router.PUT("/plants/:id/journal/:entryId", plantHandler.HandleUpdateJournalEntry)

	// 404 handler
	router.NoRoute(plantHandler.HandlePlantList) // Redirects all unknown routes to plant list

	return router, nil
}
