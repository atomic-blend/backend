package folder

import (
	"atomic_blend_api/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteFolder handles the deletion of a folder
func (c *Controller) DeleteFolder(ctx *gin.Context) {
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

	// Delete folder in database
	err = c.folderRepo.Delete(ctx, folderObjectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting folder: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
