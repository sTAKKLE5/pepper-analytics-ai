package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"pepper-analytics-ai/internal/database"
	"pepper-analytics-ai/internal/routes"
)

func loadEnv() error {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Try to load the environment-specific file first
	envFile := filepath.Join(wd, ".env."+env)
	err = godotenv.Load(envFile)
	if err != nil {
		log.Printf("Warning: Error loading %s: %v", envFile, err)

		// If environment-specific file fails, try loading default .env
		err = godotenv.Load(filepath.Join(wd, ".env"))
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// Load environment variables
	if err := loadEnv(); err != nil {
		log.Fatalf("Failed to load environment: %v", err)
	}

	// Initialize database
	dbConfig := database.NewDefaultConfig()
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Set up router with error handling
	router, err := routes.SetupRouter(routes.RouterConfig{
		DB: db,
	})
	if err != nil {
		log.Fatalf("Failed to setup router: %v", err)
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Starting server in %s mode on port %s", os.Getenv("GIN_MODE"), port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
