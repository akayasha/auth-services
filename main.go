package main

import (
	config "auth-services/config"
	"auth-services/routes"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Warning: .env file not found, using default values")
	}

	// Initialize database connection
	config.ConnectDatabase()

	// Setup routes
	r := routes.SetupRouter()

	// Start the server on port 8081
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
