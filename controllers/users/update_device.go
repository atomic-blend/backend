package users

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// UpdateDeviceRequest represents the data needed to update a user device
type UpdateDeviceRequest struct {
	DeviceID   string `json:"deviceId" binding:"required"`
	DeviceName string `json:"deviceName" binding:"required"`
	FcmToken   string `json:"fcmToken" binding:"required"`
}

// UpdateDeviceInfo allows users to add or update device information
func (c *UserController) UpdateDeviceInfo(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Parse and validate device update request
	var deviceReq UpdateDeviceRequest
	if err := ctx.ShouldBindJSON(&deviceReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Fetch the current user data
	user, err := c.userRepo.FindByID(ctx, authUser.UserID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve user profile for device update")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	// Create device object from request
	newDevice := &models.UserDevice{
		DeviceID:   deviceReq.DeviceID,
		DeviceName: deviceReq.DeviceName,
		FcmToken:   deviceReq.FcmToken,
	}

	// Check if devices array is initialized
	if user.Devices == nil {
		user.Devices = make([]*models.UserDevice, 0)
	}

	// Check if the device already exists, if so update it
	deviceFound := false
	for i, device := range user.Devices {
		if device.DeviceID == deviceReq.DeviceID {
			// Update existing device
			user.Devices[i] = newDevice
			deviceFound = true
			break
		}
	}

	// If device not found, add it
	if !deviceFound {
		user.Devices = append(user.Devices, newDevice)
	}

	// Update user in database
	updatedUser, err := c.userRepo.Update(ctx, user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update user device information")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update device information"})
		return
	}

	// Remove sensitive data before sending response
	updatedUser.Password = nil

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Device information updated successfully",
		"data":    updatedUser,
	})
}
