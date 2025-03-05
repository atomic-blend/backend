package auth

import (
	"atomic_blend_api/utils/jwt"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/rs/zerolog/log"
)

// GetJWKS provides the JSON Web Key Set for JWT validation
// @Summary Get JWKS (JSON Web Key Set)
// @Description Returns the JWKS that can be used to verify JWT tokens issued by this server
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /auth/.well-known/jwks.json [get]
func (c *Controller) GetJWKS(ctx *gin.Context) {
	// Generate the JWK from the secret key
	secretKey := os.Getenv("SSO_SECRET")
	if secretKey == "" {
		log.Error().Msg("SSO_SECRET is not set")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Missing SSO_SECRET"})
		return
	}

	key, err := jwt.GenerateJWKS(secretKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWKS"})
		return
	}

	// Set required JWK parameters
	err = (*key).Set(jwk.KeyIDKey, "atomic-blend-sso")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set key ID"})
		return
	}

	err = (*key).Set(jwk.AlgorithmKey, "HS256")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set algorithm"})
		return
	}

	err = (*key).Set(jwk.KeyUsageKey, "sig")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set key usage"})
		return
	}

	// Create a JWKS (JSON Web Key Set) with our key
	set := jwk.NewSet()
	set.Add(*key)

	// Convert to JSON and return
	jwksJSON, err := json.Marshal(set)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JWKS"})
		return
	}

	// Set appropriate content type and return the JWKS
	ctx.Header("Content-Type", "application/json")
	ctx.Data(http.StatusOK, "application/json", jwksJSON)
}
