package controllers

import (
	"QResume/contracts"
	"QResume/service"
	"github.com/gin-gonic/gin"
	"net/http"
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
	c.JSON(http.StatusOK, gin.H{"message": "User signed on successfully"})
}

func (u *UserController) UpdateDetails(c *gin.Context) {
	// Define a struct to bind incoming JSON data
	var userDetails contracts.UserDetails

	// Bind JSON payload to the input struct
	if err := c.ShouldBindJSON(&userDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to register the user
	err := u.UserService.UpdateDetails(&userDetails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "User signed on successfully"})
}

// GetUserDetails fetches user details based on user-id
func (u *UserController) GetUserDetails(c *gin.Context) {
	userEmail := c.Param("user-email")

	// Call the service to get user details
	userDetails, err := u.UserService.GetUserDetails(userEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userDetails)
}

// GetUserQRCode fetches the user's QR code based on user-id
func (u *UserController) GetUserQRCode(c *gin.Context) {
	userEmail := c.Param("user-email")
	// Fetch the QR code
	qrCode, err := u.UserService.GetUserQRCode(userEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.File(qrCode) // Respond with the QR code file
}
