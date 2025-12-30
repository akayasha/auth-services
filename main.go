package main

import (
	"auth-services/config"
	"auth-services/routes"
	"auth-services/services"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env not found, using system env")
	}

	config.ConnectDatabase()
	config.ConnectRedis()

	// ✅ Initialize audit worker
	services.InitAuditWorker()

	r := routes.SetupRouter()

	port := config.GetEnv("APP_PORT", "8081")
	log.Printf("🚀 Server starting on port %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
