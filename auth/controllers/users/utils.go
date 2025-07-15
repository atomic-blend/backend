package users

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeletePersonalData handles the deletion of all user personal data
// This includes tasks and any other personal data associated with the user
func (c *UserController) DeletePersonalData(ctx *gin.Context, userID primitive.ObjectID) error {
	//TODO: Implement the logic to delete all personal data associated with the user with gRPC
	// For now, we will just return nil to indicate success
	return nil
}
