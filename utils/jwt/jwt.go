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
	secretKey = os.Getenv("SSO_SECRET")
	if secretKey == "" {
		secretKey = "default_access_secret" // Fallback for development
	}
	expTime = 15 * time.Minute // 15 minutes for access token
	td.ExpiresAt = time.Now().Add(expTime)

	claims := jwt.MapClaims{
		"sub":  td.UserID,
		"user_id": td.UserID,
		"aud":  "atomic-blend",
		"iss":  "atomic-blend",
		"type": string(tokenType),
		"iat":  time.Now().Unix(),
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
func ValidateToken(tokenString string, tokenType TokenType) (*jwt.MapClaims, error) {

	secretKey := os.Getenv("SSO_SECRET")
	if secretKey == "" {
		log.Error().Msg("SSO_SECRET not set")
		return nil, errors.New("SSO_SECRET not set")
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
