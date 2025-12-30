package config

import (
	"auth-services/models"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	host := GetEnv("DB_HOST", "127.0.0.1")
	user := GetEnv("DB_USER", "postgres")
	password := GetEnv("DB_PASSWORD", "")
	dbName := GetEnv("DB_NAME", "auth_service")
	port := GetEnvInt("DB_PORT", 5432)
	timezone := GetEnv("DB_TIMEZONE", "Asia/Jakarta")

	defaultDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=%d sslmode=disable TimeZone=%s",
		host, user, password, port, timezone,
	)

	sqlDB, err := sql.Open("postgres", defaultDSN)
	if err != nil {
		log.Fatalf("❌ PostgreSQL bootstrap failed: %v", err)
	}
	defer sqlDB.Close()

	var exists bool
	check := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname='%s')`, dbName)
	if err := sqlDB.QueryRow(check).Scan(&exists); err != nil {
		log.Fatalf("❌ DB existence check failed: %v", err)
	}

	if !exists {
		log.Printf("⚠️ Creating database %s", dbName)
		if _, err := sqlDB.Exec("CREATE DATABASE " + dbName); err != nil {
			log.Fatalf("❌ Create DB failed: %v", err)
		}
	}

	appDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		host, user, password, dbName, port, timezone,
	)

	DB, err = gorm.Open(postgres.Open(appDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ PostgreSQL connect failed: %v", err)
	}

	if err := DB.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.AuditLog{},
	); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	log.Println("🚀 PostgreSQL connected")
}
