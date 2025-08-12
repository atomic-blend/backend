package folder

import (
	"net/http"

	"github.com/atomic-blend/backend/productivity/models"
	"github.com/atomic-blend/backend/shared/middlewares/auth"

	"github.com/gin-gonic/gin"
)

// CreateFolder creates a new folder
// @Summary Create folder
// @Description Create a new folder
// @Tags Folders
// @Accept json
// @Produce json
// @Param folder body models.Folder true "Folder object"
// @Success 201 {object} models.Folder
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /folders [post]
func (c *Controller) CreateFolder(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var folder models.Folder
	if err := ctx.ShouldBindJSON(&folder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the user ID
	folder.UserID = authUser.UserID

	// Create folder in database
	createdFolder, err := c.folderRepo.Create(ctx, &folder)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating folder: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdFolder)
}
