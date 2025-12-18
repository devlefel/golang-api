package main

import (
	_ "device-api/docs" // Import generated docs
	"device-api/internal/domain"
	"device-api/internal/handler"
	"device-api/internal/repository"
	"device-api/internal/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Device API
// @version 1.0
// @description REST API for managing device resources.
// @host localhost:8080
// @BasePath /api/v1
func main() {
    // Load .env file if present
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

	dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = "host=db user=postgres password=postgres dbname=devices port=5432 sslmode=disable"
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Migrate the schema
    db.AutoMigrate(&domain.Device{})

    repo := repository.NewPostgresRepository(db)
    svc := service.NewDeviceService(repo)
    h := handler.NewDeviceHandler(svc)

    r := gin.Default()
    handler.RegisterRoutes(r, h)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    r.Run(":" + port)
}
