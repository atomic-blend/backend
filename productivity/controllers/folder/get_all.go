package folder

import (
	"net/http"
	"github.com/atomic-blend/backend/shared/middlewares/auth"


	"github.com/gin-gonic/gin"
)

// GetAllFolders handles the retrieval of all folders for the authenticated user
func (c *Controller) GetAllFolders(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Get folders from database
	folders, err := c.folderRepo.GetAll(ctx, authUser.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching folders: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, folders)
}
