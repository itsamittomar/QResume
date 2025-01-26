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
	// Hash the password

	// Create a new detail object
	updates := &models.Details{
		Email:         userDetails.Email,
		Linkedin:      userDetails.Linkedin,
		Github:        userDetails.Github,
		Leetcode:      userDetails.Leetcode,
		GeeksForGeeks: userDetails.GeeksForGeeks,
		Scaler:        userDetails.Scaler,
	}

	qrCodeLinks, err := u.generateQRCode(userDetails)
	if err != nil {
		fmt.Println("Error generating QR codes:", err)
		return err
	}

	// Update the QR code links in the user details object
	updates.QRCodeLeetcode = qrCodeLinks["leetcode"]
	updates.QRCodeScaler = qrCodeLinks["scaler"]
	updates.QRCodeGeeksForGeeks = qrCodeLinks["geeksforgeeks"]

	// Call the UpdateByEmail function
	err = u.UserRepo.UpdateByEmail(userDetails.Email, updates)
	if err != nil {
		fmt.Println("Error updating user details:", err)
		return err
	}

	return nil
}

func (u *UserService) generateQRCode(userDetails *contracts.UserDetails) (map[string]string, error) {
	// Map to hold the generated QR code file paths for each link
	qrCodeLinks := make(map[string]string)

	// Base directory to store QR codes
	baseDir := "qrcodes"

	// Ensure the directory exists
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		return nil, err
	}

	// Generate a QR code for each link
	links := map[string]string{
		"leetcode":      userDetails.Leetcode,
		"scaler":        userDetails.Scaler,
		"geeksforgeeks": userDetails.GeeksForGeeks,
	}

	for platform, url := range links {
		if url != "" { // Skip if URL is empty
			// Define the filename for storing the QR code
			fileName := fmt.Sprintf("%s/%s_%s.png", baseDir, userDetails.Email, platform)

			// Generate the QR code
			err := qrcode.WriteFile(url, qrcode.Medium, 256, fileName)
			if err != nil {
				return nil, fmt.Errorf("failed to generate QR code for %s: %w", platform, err)
			}

			// Map the platform to its QR code file path or URL
			qrCodeLinks[platform] = fmt.Sprintf("http://localhost:8080/static/%s", fileName)
		}
	}

	return qrCodeLinks, nil
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
