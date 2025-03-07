package user_role

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteRole deletes a user role
// @Summary Delete user role
// @Description Delete a user role by ID
// @Tags User Roles
// @Param id path string true "Role ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/user-roles/{id} [delete]
func (c *UserRoleController) DeleteRole(ctx *gin.Context) {
	id := ctx.Param("id")
	
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	// Delete the role
	err = c.userRoleRepo.Delete(ctx, objID)
	if err != nil {
		if err.Error() == "user role not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User role not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.Status(http.StatusNoContent)
}