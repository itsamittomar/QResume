package main

import (
	"QResume/controllers"
	"QResume/repo"
	"QResume/service"
	"QResume/models"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Database connection
	dsn := "user:password@tcp(db:3306)/qresume?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Automigrate the User model
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}

	// Inject dependencies
	userRepo := repo.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Initialize Gin router
	r := gin.Default()

	// Define routes
	r.POST("/register", userController.RegisterUser)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
