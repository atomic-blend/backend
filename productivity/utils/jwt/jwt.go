package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TokenType represents the type of token
type TokenType string

const (
	// AccessToken is used for authenticating requests
	AccessToken TokenType = "access"

	// RefreshToken is used to get a new access token
	RefreshToken TokenType = "refresh"
)

// TokenDetails contains the token information
type TokenDetails struct {
	Token     string
	TokenType TokenType
	ExpiresAt time.Time
	UserID    string
}

// CustomClaims represents the custom claims in the JWT
type CustomClaims struct {
	UserID       *string `json:"user_id"`
	IsSubscribed *bool   `json:"is_subscribed"`
	Type         *string `json:"type"`
	Roles        *[]string `json:"roles"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token
func GenerateToken(userID primitive.ObjectID, tokenType TokenType) (*TokenDetails, error) {
	var td TokenDetails
	td.UserID = userID.Hex()
	td.TokenType = tokenType

	var secretKey string
	var expTime time.Duration

	// Set different expiry times and secrets based on token type
	secretKey = os.Getenv("SSO_SECRET")
	if secretKey == "" {
		secretKey = "default_access_secret" // Fallback for development
	}
	expTime = 15 * time.Minute // 15 minutes for access token
	td.ExpiresAt = time.Now().Add(expTime)

	claims := jwt.MapClaims{
		"sub":     td.UserID,
		"user_id": td.UserID,
		"aud":     "atomic-blend",
		"iss":     "atomic-blend",
		"type":    string(tokenType),
		"iat":     time.Now().Unix(),
	}

	if tokenType == AccessToken {
		claims["exp"] = td.ExpiresAt.Unix()
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
func ValidateToken(tokenString string, tokenType TokenType) (*CustomClaims, error) {
	secretKey := os.Getenv("SSO_SECRET")
	if secretKey == "" {
		log.Error().Msg("SSO_SECRET not set")
		return nil, errors.New("SSO_SECRET not set")
	}

	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Verify token type
	if *claims.Type != string(tokenType) {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

// GenerateJWKS creates a new JSON Web Key Set
func GenerateJWKS(secretType string) (*jwk.Key, error) {
	secretKey := os.Getenv("SSO_SECRET")
	if secretKey == "" {
		return nil, errors.New("SSO_SECRET not set")
	}

	log.Debug().Msgf("Generating JWK with secret: %s\n", secretKey)

	key, err := jwk.New([]byte(secretKey))
	log.Debug().Msgf("Generated JWK: %v\n", err)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWK: %s", err)
	}
	if _, ok := key.(jwk.SymmetricKey); !ok {
		log.Error().Msgf("expected jwk.SymmetricKey, got %T\n", key)
		return nil, errors.New("failed to create JWK")
	}

	return &key, nil
}
