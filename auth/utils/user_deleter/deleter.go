// Package userdeleter provides utilities for completely deleting a user
package userdeleter

import (
	"context"

	"connectrpc.com/connect"
	authv1 "github.com/atomic-blend/backend/grpc/gen/auth/v1"
	productivityv1 "github.com/atomic-blend/backend/grpc/gen/productivity/v1"
	productivityclient "github.com/atomic-blend/backend/shared/grpc/productivity"
	"github.com/atomic-blend/backend/shared/repositories/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeletePersonalDataAndUser handles the deletion of all user personal data
// This includes tasks and any other personal data associated with the user
func DeletePersonalDataAndUser(userID primitive.ObjectID, productivityClient productivityclient.Interface, userRepo user.Interface) error {
	// Create the request with the correct Connect-RPC format
	req := connect.NewRequest(&productivityv1.DeleteUserDataRequest{
		User: &authv1.User{
			Id: userID.Hex(),
		},
	})

	// Call the service with the appropriate context
	_, err := productivityClient.DeleteUserData(context.TODO(), req)
	if err != nil {
		return err
	}

	// Delete the user
	err = userRepo.Delete(context.TODO(), userID.Hex())
	if err != nil {
		return err
	}

	return nil
}
