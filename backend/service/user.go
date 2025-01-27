package service

import (
	"QResume/contracts"
	"QResume/models"
	"QResume/repo"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/skip2/go-qrcode"
	"golang.org/x/crypto/bcrypt"
	"net/url"
	"os"
	"path/filepath"
)

type UserService struct {
	UserRepo *repo.UserRepo
}

// NewUserService initializes a new UserService
func NewUserService(userRepo *repo.UserRepo) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

// RegisterUser handles the business logic for user registration
func (u *UserService) RegisterUser(userDetails *contracts.Register) error {
	// Hash the password
	hashedPassword, err := u.HashPassword(userDetails.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create a new user object
	user := &models.User{
		Email:    userDetails.Email,
		Password: hashedPassword,
	}

	// Save the user via the repository
	err = u.UserRepo.Register(user)
	if err != nil {
		// Handle duplicate entry error
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			err := u.Login(userDetails.Email, userDetails.Password)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("failed to register user: %w", err)
		}
	}

	return nil
}

func (u *UserService) UpdateDetails(userDetails *contracts.UserDetails) error {
	// Create a new User object with the updated details
	updates := &models.User{
		Email: userDetails.Email,
		Name:  userDetails.Name,
		Details: models.Details{
			Email:         userDetails.Email,
			Phone:         userDetails.Phone,
			Linkedin:      userDetails.Linkedin,
			Github:        userDetails.Github,
			Leetcode:      userDetails.Leetcode,
			GeeksForGeeks: userDetails.GeeksForGeeks,
			Scaler:        userDetails.Scaler,
		},
	}

	// Generate the QR Code and store it in the User object (assuming QR code is part of User)
	qrCodeURL, err := u.generateQRCode(userDetails)
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return err
	}
	updates.QR = []byte(qrCodeURL) // Assuming the QR code is stored in the QR field

	// Call the UpdateByEmail function to update the details in the repository
	err = u.UserRepo.UpdateByEmail(userDetails.Email, updates)
	if err != nil {
		fmt.Println("Error updating user details:", err)
		return err
	}

	return nil
}

func (u *UserService) generateQRCode(userDetails *contracts.UserDetails) (string, error) {
	// Define the directory for storing QR codes
	// This directory should be volume-mapped in your Docker container
	dir := "/app/qrcodes"

	// Get the absolute directory path (useful for logging)
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute directory path: %w", err)
	}
	fmt.Println("Using directory path:", absDir)

	// Check if the directory exists, and create it if it doesn't
	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		fmt.Println("Directory does not exist. Attempting to create...")
		if err := os.MkdirAll(absDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create QR code directory: %w", err)
		}
		fmt.Println("Directory created successfully at:", absDir)
	} else {
		fmt.Println("Directory already exists at:", absDir)
	}

	// Combine all links into a single string
	combinedURL := fmt.Sprintf(
		"http://localhost:8080/user?leetcode=%s&scaler=%s&geeksforgeeks=%s",
		url.QueryEscape(userDetails.Leetcode),
		url.QueryEscape(userDetails.Scaler),
		url.QueryEscape(userDetails.GeeksForGeeks),
	)

	// Define the filename for the QR code (using filepath for cross-platform compatibility)
	fileName := filepath.Join(absDir, fmt.Sprintf("%s_combined.png", userDetails.Email))

	// Generate and save the QR code to the file
	if err := qrcode.WriteFile(combinedURL, qrcode.Medium, 256, fileName); err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}
	fmt.Println("QR code successfully generated and saved at:", fileName)

	// Return the path to the saved QR code
	return fileName, nil
}

// Login handles user login by verifying the password
func (u *UserService) Login(email, password string) error {
	// Retrieve user details from the repository
	user, err := u.UserRepo.Login(email) // Assuming GetUserByEmail fetches user by email
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify the password
	err = u.CheckPassword(user.Password, password)
	if err != nil {
		return fmt.Errorf("invalid credentials: %w", err)
	}

	return nil
}

// HashPassword hashes a plain-text password
func (u *UserService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword compares a hashed password with a plain-text password
func (u *UserService) CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GetUserDetails fetches user details by ID
func (u *UserService) GetUserDetails(userID string) (*models.User, error) {
	return u.UserRepo.GetUserByID(userID)
}

// GetUserQRCode fetches the path of the user's QR code by ID
func (u *UserService) GetUserQRCode(userID string) (string, error) {
	user, err := u.UserRepo.GetUserByID(userID)
	if err != nil {
		return "", err
	}

	// Check if the QR code exists
	qrCodePath := string(user.QR)
	if _, err := os.Stat(qrCodePath); os.IsNotExist(err) {
		return "", fmt.Errorf("QR code not found for user %s", userID)
	}

	return qrCodePath, nil
}