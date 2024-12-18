package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	model "github.com/liju-github/EcommerceUserService/models"
	userPb "github.com/liju-github/EcommerceUserService/proto/user"
	"github.com/liju-github/EcommerceUserService/repository"
	util "github.com/liju-github/EcommerceUserService/utils"
)

const (
	TokenExpiry = 24 * time.Hour
)

type UserService struct {
	userPb.UnimplementedUserServiceServer
	repo repository.UserRepository
}

type CustomClaims struct {
	UserID     string `json:"userId"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Reputation int32  `json:"reputation"`
	jwt.RegisteredClaims
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}
func (s *UserService) GetAllUsers(ctx context.Context, req *userPb.GetAllUsersRequest) (*userPb.GetAllUsersResponse, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	var userResponses []*userPb.User
	for _, user := range users {
		userResponses = append(userResponses, &userPb.User{
			Id:          user.ID,
			Email:       user.Email,
			Name:        user.Name,
			Reputation:  user.Reputation,
			StreetName:  user.StreetName,
			Locality:    user.Locality,
			State:       user.State,
			Pincode:     user.Pincode,
			PhoneNumber: user.PhoneNumber,
			IsVerified:  user.IsVerified,
			IsBanned:    user.IsBanned,
		})
	}

	return &userPb.GetAllUsersResponse{
		Success: true,
		Users:   userResponses,
	}, nil
}

func (s *UserService) Register(ctx context.Context, req *userPb.RegisterRequest) (*userPb.RegisterResponse, error) {
	existingUser, err := s.repo.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, model.ErrDuplicateEmail
	}

	// Generate password hash
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate verification code (you might want to use a more sophisticated method)
	verificationCode := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)

	user := model.User{
		ID:               fmt.Sprintf("usr_%d", time.Now().UnixNano()),
		Email:            req.Email,
		PasswordHash:     string(passwordHash),
		Name:             req.Name,
		StreetName:       req.StreetName,
		Locality:         req.Locality,
		State:            req.State,
		Pincode:          req.Pincode,
		PhoneNumber:      req.PhoneNumber,
		Reputation:       0,
		IsVerified:       false,
		VerificationCode: verificationCode,
	}

	if err := s.repo.CreateUser(&user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Here you would typically send an email with the verification code
	// sendVerificationEmail(user.Email, verificationCode)

	return &userPb.RegisterResponse{
		Success: true,
		Message: "Registration successful. Please check your email for verification.",
	}, nil
}

// Login verifies credentials and returns a token
func (s *UserService) Login(ctx context.Context, req *userPb.LoginRequest) (*userPb.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, model.ErrInvalidPassword
	}

	return &userPb.LoginResponse{
		Success: true,
		UserId:  user.ID,
	}, nil
}

// VerifyEmail handles email verification
func (s *UserService) VerifyEmail(ctx context.Context, req *userPb.EmailVerificationRequest) (*userPb.EmailVerificationResponse, error) {
	user, err := s.repo.GetUserByEmail(req.UserId)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	if user.VerificationCode != req.VerificationCode {
		return nil, model.ErrInvalidCode
	}

	if err := s.repo.UpdateUserVerification(user.ID, true); err != nil {
		return nil, fmt.Errorf("failed to update verification status: %w", err)
	}

	return &userPb.EmailVerificationResponse{
		Success: true,
		Message: "Email successfully verified",
	}, nil
}

// GetProfile retrieves user profile
func (s *UserService) GetProfile(ctx context.Context, req *userPb.ProfileRequest) (*userPb.ProfileResponse, error) {
	user, err := s.repo.GetUserProfile(req.UserId)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	return &userPb.ProfileResponse{
		UserId:      user.ID,
		Email:       user.Email,
		Name:        user.Name,
		Reputation:  user.Reputation,
		StreetName:  user.StreetName,
		Locality:    user.Locality,
		State:       user.State,
		Pincode:     user.Pincode,
		PhoneNumber: user.PhoneNumber,
		IsVerified:  user.IsVerified,
	}, nil
}

// VerifyTokenMiddleware middleware for token verification
func (s *UserService) GetUserByToken(ctx context.Context, req *userPb.GetUserByTokenRequest) (*userPb.ProfileResponse, error) {
	claims, err := util.ValidateToken(req.Token)
	if err != nil {
		return nil, model.ErrInvalidToken
	}

	user, err := s.repo.GetUserByID(claims.UserID)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	fmt.Println("user record is ", user)

	if !user.IsVerified {
		return nil, model.ErrUserNotVerified
	}
	response := &userPb.ProfileResponse{
		UserId:      user.ID,
		Email:       user.Email,
		Name:        user.Name,
		Reputation:  user.Reputation,
		StreetName:  user.StreetName,
		Locality:    user.Locality,
		State:       user.State,
		Pincode:     user.Pincode,
		PhoneNumber: user.PhoneNumber,
		IsVerified:  user.IsVerified,
		IsBanned:    false,
	}

	return response, nil
}

// UpdateProfile updates user profile information
func (s *UserService) UpdateProfile(ctx context.Context, req *userPb.UpdateProfileRequest) (*userPb.UpdateProfileResponse, error) {
	// Verify token first
	fmt.Println("the request is ", req)

	// Fetch user by ID
	user, err := s.repo.GetUserByID(req.UserId)
	if err != nil {
		return nil, err
	}
	fmt.Println("the profile is ", req)

	// Update user fields if new values are provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.StreetName != "" {
		user.StreetName = req.StreetName
	}
	if req.Locality != "" {
		user.Locality = req.Locality
	}
	if req.State != "" {
		user.State = req.State
	}
	if req.Pincode != "" {
		user.Pincode = req.Pincode
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}

	// Save updated user profile in repository
	if err := s.repo.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// Refetch the updated user data to ensure data consistency
	user, err = s.repo.GetUserByID(req.UserId)
	if err != nil {
		return nil, err
	}

	// Prepare response with updated profile details
	return &userPb.UpdateProfileResponse{
		Success: true,
		Message: "Profile updated successfullyyy",
		Profile: &userPb.ProfileResponse{
			UserId:      user.ID,
			Email:       user.Email,
			Name:        user.Name,
			Reputation:  user.Reputation,
			StreetName:  user.StreetName,
			Locality:    user.Locality,
			State:       user.State,
			Pincode:     user.Pincode,
			PhoneNumber: user.PhoneNumber,
			IsVerified:  user.IsVerified,
			IsBanned:    user.IsBanned,
		},
	}, nil
}

func (s *UserService) CheckBan(ctx context.Context, req *userPb.CheckBanRequest) (*userPb.CheckBanResponse, error) {

	status, error := s.repo.CheckBan(req.UserID)

	return &userPb.CheckBanResponse{
		BanStatus: status,
	}, error

}

func (s *UserService) BanUser(ctx context.Context, req *userPb.BanUserRequest) (*userPb.BanUserResponse, error) {
	if req.UserId == "" {
		return &userPb.BanUserResponse{
			Success: false,
			Message: "User Ban failed",
		}, errors.New("userId doesnt exist")
	}

	if err := s.repo.BanUser(req.UserId); err != nil {
		return &userPb.BanUserResponse{
			Success: false,
			Message: "User Ban failed",
		}, errors.New(err.Error())
	}

	return &userPb.BanUserResponse{
		Success: true,
		Message: "User Banned Succesfully",
	}, nil

}

func (s *UserService) UnBanUser(ctx context.Context, req *userPb.UnBanUserRequest) (*userPb.UnBanUserResponse, error) {
	if req.UserId == "" {
		return &userPb.UnBanUserResponse{
			Success: false,
			Message: "User UnBan failed",
		}, errors.New("userId doesnt exist")
	}

	if err := s.repo.UnBanUser(req.UserId); err != nil {
		return &userPb.UnBanUserResponse{
			Success: false,
			Message: "User UnBan failed",
		}, errors.New(err.Error())
	}

	return &userPb.UnBanUserResponse{
		Success: true,
		Message: "User UnBanned Succesfully",
	}, nil
}
