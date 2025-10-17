// Package utils provides utility functions for the auth service.
package utils

import (
	"os"
	"strconv"

	"github.com/atomic-blend/backend/auth/repositories"
	"github.com/atomic-blend/backend/shared/repositories/user"
	"github.com/gin-gonic/gin"
)

// GetRemainingSpots returns the remaining spots for the auth service.
func GetRemainingSpots(ctx *gin.Context, repository user.Interface) (int64, error) {
	maxUsersString := os.Getenv("AUTH_MAX_NB_USER")
	maxUsers := int64(1)
	if maxUsersString != "" {
		maxUsersInt, err := strconv.ParseInt(maxUsersString, 10, 64)
		if err != nil {
			return 0, err
		}
		maxUsers = maxUsersInt
	}

	users, err := repository.Count(ctx.Request.Context())
	if err != nil {
		return 0, err
	}
	currentUserCount := users

	return maxUsers - currentUserCount, nil
}

// GetUserCode returns the code for a user in the waiting list by email
func GetUserCode(ctx *gin.Context, waitingListRepo repositories.WaitingListRepositoryInterface, email string) (*string, error) {
	record, err := waitingListRepo.GetByEmail(ctx.Request.Context(), email)
	if err != nil {
		return nil, err
	}
	return record.Code, nil
}
