package service

import (
	"QResume/contracts"
	"QResume/models"
	"QResume/repo"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql" // Import MySQL driver for error handling
	"golang.org/x/crypto/bcrypt"
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
			if err != nil{
				return err
			}
		} else{
			return fmt.Errorf("failed to register user: %w", err)
		}
	}

	return nil
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
