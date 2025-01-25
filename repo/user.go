package repo

import (
	"QResume/models"
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
