package config

import (
	"auth-services/models"
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var DB *gorm.DB

func ConnectDatabase() {
	host := getEnv("DB_HOST", "127.0.0.1")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "pass")
	dbName := getEnv("DB_NAME", "go_user")
	port := getEnvAsInt("DB_PORT", 5432)
	timezone := getEnv("DB_TIMEZONE", "Asia/Jakarta")

	// Step 1: Connect to default DB
	defaultDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=%d sslmode=disable TimeZone=%s",
		host, user, password, port, timezone,
	)

	sqlDB, err := sql.Open("postgres", defaultDSN)
	if err != nil {
		log.Fatalf("❌ Failed connecting to default PostgreSQL DB: %v", err)
	}
	defer sqlDB.Close()

	// Step 2: Check if DB exists
	var exists bool
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname='%s')", dbName)

	err = sqlDB.QueryRow(checkQuery).Scan(&exists)
	if err != nil {
		log.Fatalf("❌ Failed checking database existence: %v", err)
	}

	// Step 3: Create DB if not exists
	if !exists {
		log.Printf("⚠️ Database %s does not exist — creating...", dbName)
		_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			log.Fatalf("❌ Failed to create database: %v", err)
		}
		log.Println("✅ Database created successfully!")
	}

	// Step 4: Connect to the actual application DB
	appDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		host, user, password, dbName, port, timezone,
	)

	DB, err = gorm.Open(postgres.Open(appDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL (%s): %v", dbName, err)
	}

	// Step 5: Run migrations
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	log.Println("🚀 PostgreSQL connected & ready (DB:", dbName, ")")
}

// getEnv retrieves the value of the environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves the value of the environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
