package controllers

import (
	"net/http"
	"QResume/service"
	"QResume/contracts"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *service.UserService
}

// NewUserController initializes a new UserController
func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

// RegisterUser handles the registration of a new user
func (u *UserController) RegisterUser(c *gin.Context) {
	// Define a struct to bind incoming JSON data
	var userDetails contracts.Register

	// Bind JSON payload to the input struct
	if err := c.ShouldBindJSON(&userDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to register the user
	err := u.UserService.RegisterUser(&userDetails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
