package main

import (
	config "auth-services/config"
	"auth-services/routes"
	"log"
)

func main() {
	// Initialize database connection
	config.ConnectDatabase()

	// Setup routes
	r := routes.SetupRouter()

	// Start the server on port 8081
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
