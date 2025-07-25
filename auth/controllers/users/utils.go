package users

import (
	"connectrpc.com/connect"
	"github.com/atomic-blend/backend/grpc/gen/auth/v1"
	"github.com/atomic-blend/backend/grpc/gen/productivity/v1"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeletePersonalData handles the deletion of all user personal data
// This includes tasks and any other personal data associated with the user
func (c *UserController) DeletePersonalData(ctx *gin.Context, userID primitive.ObjectID) error {
	// Create the request with the correct Connect-RPC format
	req := connect.NewRequest(&productivityv1.DeleteUserDataRequest{
		User: &authv1.User{
			Id: userID.Hex(),
		},
	})

	// Call the service with the appropriate context
	_, err := c.productivityClient.DeleteUserData(ctx.Request.Context(), req)
	if err != nil {
		return err
	}

	return nil
}
