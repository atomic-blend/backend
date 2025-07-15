package folder

import (
	"net/http"
	"productivity/auth"
	"productivity/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateFolder handles the update of a folder
func (c *Controller) UpdateFolder(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	folderID := ctx.Param("id")
	if folderID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Folder ID is required"})
		return
	}

	// Convert folderID to ObjectID
	folderObjectID, err := primitive.ObjectIDFromHex(folderID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	var folder models.Folder
	if err := ctx.ShouldBindJSON(&folder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folder.UserID = authUser.UserID

	// Update folder in database
	updatedFolder, err := c.folderRepo.Update(ctx, folderObjectID, &folder)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating folder: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedFolder)
}
