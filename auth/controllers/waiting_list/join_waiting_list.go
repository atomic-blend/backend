package waitinglist

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"text/template"
	"time"

	waitinglist "github.com/atomic-blend/backend/auth/models/waiting_list"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JoinWaitingListRequest represents the structure for join waiting list request data
type JoinWaitingListRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// generateSecurityToken generates a 32-character random security token
func generateSecurityToken() (string, error) {
	bytes := make([]byte, 16) // 16 bytes = 32 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// JoinWaitingList creates a new waiting list record and returns a success message
func (c *Controller) JoinWaitingList(ctx *gin.Context) {
	var req JoinWaitingListRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if the email is already in the waiting list
	waitingListRecord, err := c.waitingListRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error_getting_waiting_list_record"})
		return
	}

	if waitingListRecord != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error_email_already_in_waiting_list"})
		return
	}

	// generate security token
	securityToken, err := generateSecurityToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error_generating_security_token"})
		return
	}

	// create a new waiting list record
	now := primitive.NewDateTimeFromTime(time.Now())
	waitingListRecord, err = c.waitingListRepo.Create(ctx, &waitinglist.WaitingList{
		Email:         req.Email,
		SecurityToken: &securityToken,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error_creating_waiting_list_record"})
		return
	}

	// get the domain from the environment variables
	domain := os.Getenv("PUBLIC_ADDRESS")
	if domain == "" {
		domain = "mail.atomic-blend.com"
	} else {
		domain = "mail." + domain
	}
	log.Info().Msgf("Domain: %s", domain)

	// template the html with gotemplate
	htmlTemplate, err := template.ParseFiles("./email_templates/join_waiting_list/join_waiting_list.html")
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse HTML template")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse HTML template"})
		return
	}

	textTemplate, err := template.ParseFiles("./email_templates/join_waiting_list/join_waiting_list.txt")
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse text template")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse text template"})
		return
	}

	// template the plain text with gotemplate
	var htmlContent bytes.Buffer
	err = htmlTemplate.Execute(&htmlContent, map[string]string{
		"email":         req.Email,
		"domain":        domain,
		"securityToken": securityToken,
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to execute HTML template")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute HTML template"})
		return
	}

	var textContent bytes.Buffer
	err = textTemplate.Execute(&textContent, map[string]string{
		"email":         req.Email,
		"domain":        domain,
		"securityToken": securityToken,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute text template")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute text template"})
		return
	}

	// Create RawMail structure for AMQP
	rawMail := map[string]interface{}{
		"headers": map[string]interface{}{
			"To":      []string{req.Email},
			"From":    "noreply@atomic-blend.com",
			"Subject": "You just joined the waiting list!",
		},
		"htmlContent":    htmlContent.String(),
		"textContent":    textContent.String(),
		"rejected":       false,
		"rewriteSubject": false,
		"graylisted":     false,
	}

	// Publish message to AMQP queue
	c.amqpService.PublishMessage("mail", "sent", map[string]interface{}{
		"waiting_list_email": true,
		"content":            rawMail,
	}, nil)

	log.Info().Msg("Waiting list email queued for sending")

	// get the position of this record (0-based index)
	position, err := c.waitingListRepo.GetPositionByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error_getting_waiting_list_position"})
		return
	}

	// get the total count after creating the record
	totalCount, err := c.waitingListRepo.Count(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error_getting_waiting_list_count"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"entry":    waitingListRecord,
		"position": position,
		"total":    totalCount,
	})
}
