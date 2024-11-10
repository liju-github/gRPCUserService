package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	model "github.com/liju-github/EcommerceUserService/models"
)

// CustomClaims extends jwt.StandardClaims
type CustomClaims struct {
	UserID     string `json:"userId"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Reputation int32  `json:"reputation"`
	jwt.RegisteredClaims
}

// JWT related constants and variables
var (
	TokenExpiry  = 24 * time.Hour
	JWTSecretKey string
)

// SetJWTSecretKey sets the JWT secret key
func SetJWTSecretKey(secret string) {
	JWTSecretKey = secret
}

// GenerateToken creates a new JWT token
func GenerateToken(user *model.User) (string, error) {
	claims := CustomClaims{
		UserID:     user.ID,
		Email:      user.Email,
		Role:       "user", // Default role
		Reputation: user.Reputation,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ecommerce-user-service",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(JWTSecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken verifies the JWT token
func ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWTSecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, model.ErrInvalidToken
	}

	return claims, nil
}
