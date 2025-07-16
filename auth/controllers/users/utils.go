package users

import (
	"net/http"

	"connectrpc.com/connect"
	grpcclientutils "github.com/atomic-blend/backend/auth/utils/grpc_client_utils"
	"github.com/atomic-blend/backend/grpc/gen/auth"
	"github.com/atomic-blend/backend/grpc/gen/productivity"
	"github.com/atomic-blend/backend/grpc/gen/productivity/productivityconnect"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeletePersonalData handles the deletion of all user personal data
// This includes tasks and any other personal data associated with the user
func (c *UserController) DeletePersonalData(ctx *gin.Context, userID primitive.ObjectID) error {
	// Create a new HTTP client
	httpClient := &http.Client{}

	// URL of the productivity service
	baseURL, err := grpcclientutils.GetServiceBaseURL("productivity")
	if err != nil {
		return err
	}

	// Create the Connect-RPC client
	client := productivityconnect.NewProductivityServiceClient(httpClient, baseURL)

	// Create the request with the correct Connect-RPC format
	req := connect.NewRequest(&productivity.DeleteUserDataRequest{
		User: &auth.User{
			Id: userID.Hex(),
		},
	})

	// Call the service with the appropriate context
	_, err = client.DeleteUserData(ctx.Request.Context(), req)
	if err != nil {
		return err
	}

	return nil
}
