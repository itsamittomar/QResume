package service

import (
	"QResume/models"
	"QResume/repo"
	"QResume/contracts"
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
func (u *UserService) RegisterUser(userDetails *contracts.Register) (*models.User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDetails.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create a new user object
	user := &models.User{
		Email:    userDetails.Email,
		Password: string(hashedPassword),
	}

	// Save the user via the repository
	err = u.UserRepo.Register(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
