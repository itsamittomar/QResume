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
	// Create a new detail object with updated user details
	updates := &models.Details{
		Email:         userDetails.Email,
		Linkedin:      userDetails.Linkedin,
		Github:        userDetails.Github,
		Leetcode:      userDetails.Leetcode,
		GeeksForGeeks: userDetails.GeeksForGeeks,
		Scaler:        userDetails.Scaler,
	}

	// Generate the combined QR Code link
	qrCodeURL, err := u.generateQRCode(userDetails)
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return err
	}

	// Update the QR Code URL in the user details (if required)
	updates.QRCodeURL = qrCodeURL

	// Call the UpdateByEmail function to update the details in the repository
	err = u.UserRepo.UpdateByEmail(userDetails.Email, updates)
	if err != nil {
		fmt.Println("Error updating user details:", err)
		return err
	}

	return nil
}

func (u *UserService) generateQRCode(userDetails *contracts.UserDetails) (string, error) {

	dir := "qrcodes"

	// Check if the directory exists, and create it if it doesn't
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("failed to create QR code directory: %w", err)
		}
	}

	// Combine all links into a single URL string
	combinedURL := fmt.Sprintf(
		"http://localhost:8080/user?leetcode=%s&scaler=%s&geeksforgeeks=%s",
		url.QueryEscape(userDetails.Leetcode),
		url.QueryEscape(userDetails.Scaler),
		url.QueryEscape(userDetails.GeeksForGeeks),
	)

	// Define the filename for storing the QR code
	fileName := fmt.Sprintf("qrcodes/%s_combined.png", userDetails.Email) // Save in the 'qrcodes' directory

	// Generate the QR code
	err := qrcode.WriteFile(combinedURL, qrcode.Medium, 256, fileName)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Return the URL to access the QR code
	return fmt.Sprintf("http://localhost:8080/static/%s", fileName), nil
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
