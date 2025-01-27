package repo

import (
	"QResume/models"
	"errors"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

// NewUserRepo initializes a new UserRepo
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

// Register saves the user to the database
func (u *UserRepo) Register(user *models.User) error {
	if err := u.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) UpdateByEmail(email string, updates *models.User) error {
	// Use GORM to update the fields based on the email
	if err := u.DB.Model(&models.User{}).Where("email = ?", email).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

// Login checks if the email and password match a user in the database
func (u *UserRepo) Login(email string) (*models.User, error) {
	var user models.User

	// Find user by email
	if err := u.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Return the user if login is successful
	return &user, nil
}

// GetUserByID fetches a user by ID
func (u *UserRepo) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := u.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) GetUserByEmail(userEmail string) (*models.User, error) {
	var user models.User
	if err := u.DB.Where("email = ?", userEmail).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
