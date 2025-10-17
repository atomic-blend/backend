// Package config provides configuration-related HTTP handlers for the auth service.
package config

import (
	"net/http"
	"os"
	"strconv"
	"strings"

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

	maxUsersString := os.Getenv("AUTH_MAX_NB_USER")
	maxUsers := int64(1)
	if maxUsersString != "" {
		maxUsersInt, err := strconv.ParseInt(maxUsersString, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AUTH_MAX_NB_USER"})
			return
		}
		maxUsers = maxUsersInt
	}

	users, err := c.userRepo.Count(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current user count"})
		return
	}
	currentUserCount := users

	ctx.JSON(http.StatusOK, gin.H{
		"domains":        domains,
		"remainingSpots": maxUsers - currentUserCount,
	})
}
