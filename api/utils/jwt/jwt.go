package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type TokenDetails struct {
	Token     string
	TokenType TokenType
	ExpiresAt time.Time
	UserID    string
}

// GenerateToken creates a new JWT token
func GenerateToken(userID primitive.ObjectID, tokenType TokenType) (*TokenDetails, error) {
	var td TokenDetails
	td.UserID = userID.Hex()
	td.TokenType = tokenType

	var secretKey string
	var expTime time.Duration

	// Set different expiry times and secrets based on token type
	if tokenType == AccessToken {
		secretKey = os.Getenv("ACCESS_TOKEN_SECRET")
		if secretKey == "" {
			secretKey = "default_access_secret" // Fallback for development
		}
		expTime = 15 * time.Minute // 15 minutes for access token
	} else {
		secretKey = os.Getenv("REFRESH_TOKEN_SECRET")
		if secretKey == "" {
			secretKey = "default_refresh_secret" // Fallback for development
		}
		expTime = 7 * 24 * time.Hour // 7 days for refresh token
	}

	td.ExpiresAt = time.Now().Add(expTime)

	claims := jwt.MapClaims{
		"user_id": td.UserID,
		"type":    string(tokenType),
		"exp":     td.ExpiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	var err error
	td.Token, err = token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &td, nil
}

// ValidateToken verifies if a token is valid
func ValidateToken(tokenString string, tokenType TokenType) (*jwt.MapClaims, error) {
	var secretKey string

	if tokenType == AccessToken {
		secretKey = os.Getenv("ACCESS_TOKEN_SECRET")
		if secretKey == "" {
			secretKey = "default_access_secret"
		}
	} else {
		secretKey = os.Getenv("REFRESH_TOKEN_SECRET")
		if secretKey == "" {
			secretKey = "default_refresh_secret"
		}
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Verify token type
	if claims["type"] != string(tokenType) {
		return nil, errors.New("invalid token type")
	}

	return &claims, nil
}
