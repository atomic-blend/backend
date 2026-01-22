// Package config provides configuration-related HTTP handlers for the auth service.
package config

import (
	"net/http"
	"os"
	"strings"

	"github.com/atomic-blend/backend/auth/utils"
	"github.com/gin-gonic/gin"
)

// GetConfig returns the available account domains
// @Summary Returns the available account domains
// @Description Returns the available account domains from the environment variable
// @Tags Config
// @Accept json
// @Produce json
// @Success 200 {array} string
// @Router /config [get]
func (c *Controller) GetConfig(ctx *gin.Context) {
	domainList := os.Getenv("ACCOUNT_DOMAINS")
	domains := []string{}
	if domainList != "" {
		domains = strings.Split(domainList, ",")
	}

	spotsRemaining, err := utils.GetRemainingSpots(ctx, c.userRepo, c.waitingListRepo)
	if err != nil {
		// Check if it's a parsing error by checking the error message
		if err.Error() == "strconv.ParseInt: parsing \"invalid\": invalid syntax" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AUTH_MAX_NB_USER"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current user count"})
		}
		return
	}

	// get payment config
	isPaymentEnabled := os.Getenv("PAYMENT_ENABLED") == "true"
	
	ctx.JSON(http.StatusOK, gin.H{
		"domains":        domains,
		"remainingSpots": spotsRemaining,
		"paymentEnabled": isPaymentEnabled,
	})
}
