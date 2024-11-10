package repository

import (
	"errors"
	"fmt"

	model "github.com/liju-github/EcommerceUserService/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserByID(id string) (*model.User, error)
	UpdateUserVerification(userID string, isVerified bool) error
	GetUserProfile(userID string) (*model.User, error)
	UpdateUser(user *model.User) error
	StoreVerificationCode(userID, code string) error
	GetVerificationCode(userID string) (string, error)
	CheckBan(userID string)(bool,error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// CreateUser creates a new user record
func (r *userRepository) CreateUser(user *model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByEmail retrieves a user by their email address
func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

// GetUserByID retrieves a user by their ID
func (r *userRepository) GetUserByID(id string) (*model.User, error) {

	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

// UpdateUserVerification updates the verification status of a user
func (r *userRepository) UpdateUserVerification(userID string, isVerified bool) error {
	result := r.db.Model(&model.User{}).Where("id = ?", userID).Update("is_verified", isVerified)
	if result.Error != nil {
		return fmt.Errorf("failed to update user verification: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// GetUserProfile retrieves the user profile by userID
func (r *userRepository) GetUserProfile(userID string) (*model.User, error) {
	var user model.User
	if err := r.db.Select("id, email, name, street_name, locality, state, pincode, phone_number, reputation, is_verified").
		Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	return &user, nil
}

// UpdateUser updates a user's information
func (r *userRepository) UpdateUser(user *model.User) error {
	result := r.db.Model(&model.User{}).Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"name":        user.Name,
			"street_name": user.StreetName,
			"locality":    user.Locality,
			"state":       user.State,
			"pincode":     user.Pincode,
			"phone_number": user.PhoneNumber,
		})
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// StoreVerificationCode stores the verification code for a user
func (r *userRepository) StoreVerificationCode(userID, code string) error {
	result := r.db.Model(&model.User{}).Where("id = ?", userID).Update("verification_code", code)
	if result.Error != nil {
		return fmt.Errorf("failed to store verification code: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// GetVerificationCode retrieves the verification code for a user
func (r *userRepository) GetVerificationCode(userID string) (string, error) {
	var user model.User
	if err := r.db.Select("verification_code").Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("user not found")
		}
		return "", fmt.Errorf("failed to get verification code: %w", err)
	}
	return user.VerificationCode, nil
}

func (r *userRepository)CheckBan(userID string) (bool,error){
	var user model.User
	if err:=r.db.Select("verification_code").Where("id = ?", userID).First(&user).Error;err!=nil{
		return true, errors.New("check ban failed")
	}
	if user.IsBanned{
		return true,nil
	}
	return false,nil
}