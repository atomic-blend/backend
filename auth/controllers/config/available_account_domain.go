// Package config provides configuration-related HTTP handlers for the auth service.
package config

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// AvailableAccountDomain returns the available account domains
// @Summary Returns the available account domains
// @Description Returns the available account domains from the environment variable
// @Tags Domain
// @Accept json
// @Produce json
// @Success 200 {array} string
// @Router /domain/available [get]
func (c *Controller) AvailableAccountDomain(ctx *gin.Context) {
	domainList := os.Getenv("ACCOUNT_DOMAINS")
	if domainList == "" {
		ctx.JSON(http.StatusOK, gin.H{"domains": []string{}})
		return
	}
	domains := strings.Split(domainList, ",")
	ctx.JSON(http.StatusOK, gin.H{"domains": domains})
}
