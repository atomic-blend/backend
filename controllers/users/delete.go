package users

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/repositories"
	"atomic_blend_api/utils/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TaskRepositoryFactory is a function type to create task repositories
type TaskRepositoryFactory func() repositories.TaskRepositoryInterface

// Default implementation of TaskRepositoryFactory that will be used in production
var defaultTaskRepositoryFactory TaskRepositoryFactory = func() repositories.TaskRepositoryInterface {
	return repositories.NewTaskRepository(db.Database)
}

// DeleteAccount handles user account deletion
// @Summary Delete user account
// @Description Permanently delete the authenticated user's account
// @Tags Users
// @Produce json
// @Success 200 {object} map[string]interface{} "Successfully deleted account"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/me [delete]
func (c *UserController) DeleteAccount(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Convert the string ID to ObjectID if needed
	userID := authUser.UserID

	// Get the user from database to confirm they exist
	user, err := c.userRepo.FindByID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user"})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete all personal data first
	if err := c.DeletePersonalData(ctx, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete personal data: " + err.Error()})
		return
	}

	// Delete the user
	err = c.userRepo.Delete(ctx, userID.Hex())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account: " + err.Error()})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Account successfully deleted",
	})
}



// getTaskRepository returns a task repository instance using the factory
func (c *UserController) getTaskRepository() repositories.TaskRepositoryInterface {
	return defaultTaskRepositoryFactory()
}
