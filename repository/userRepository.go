package repository

import (
	"auth-services/models"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(users *models.User) error
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByRole(role string) (*models.User, error)
	FindByDob(dob time.Time) (*models.User, error)
	UpdateUser(user *models.User) error
	FindByName(name string) (*models.User, error)
	FindByUUID(uuid string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) FindByUUID(uuid string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	if db == nil {
		panic("Database connection is nil")
	}
	return &userRepository{db: db}
}

// Create Users
func (r *userRepository) CreateUser(users *models.User) error {
	return r.db.Create(users).Error
}

// Find User By Username
func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Find User by Email
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Find User By Role
func (r *userRepository) FindByRole(role string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("role = ?", role).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Find User By DOB
func (r *userRepository) FindByDob(dob time.Time) (*models.User, error) {
	var user models.User
	if err := r.db.Where("dob = ?", dob).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update User
func (r *userRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

// Find By Name
func (r *userRepository) FindByName(fistName string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("name = ?", fistName).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
