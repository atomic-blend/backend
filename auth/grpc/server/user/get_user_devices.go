package server

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	userv1 "github.com/atomic-blend/backend/grpc/gen/user/v1"
)

// GetUserDevices is the gRPC method to retrieve user devices
func (userGrpcServer *UserGrpcServer) GetUserDevices(ctx context.Context, req *connect.Request[userv1.GetUserDevicesRequest]) (*connect.Response[userv1.GetUserDevicesResponse], error) {
	// Call the repository method to get user devices
	user, err := userGrpcServer.userRepo.GetByID(ctx, req.Msg.User.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get user devices: %w", err))
	}

	devices := []*userv1.UserDevice{}
	for _, device := range user.Devices {
		deviceTz := ""
		if device.DeviceTimezone != nil {
			deviceTz = *device.DeviceTimezone
		}
		devices = append(devices, &userv1.UserDevice{
			DeviceId:       device.DeviceID,
			DeviceName:     device.DeviceName,
			FcmToken:       device.FcmToken,
			DeviceTimezone: &deviceTz,
		})	
	}

	// Create the response
	resp := &userv1.GetUserDevicesResponse{
		Devices: devices,
	}

	return connect.NewResponse(resp), nil
}

// GetUserPublicKey is the gRPC method to retrieve user public key
func (userGrpcServer *UserGrpcServer) GetUserPublicKey(ctx context.Context, req *connect.Request[userv1.GetUserPublicKeyRequest]) (*connect.Response[userv1.GetUserPublicKeyResponse], error) {
	// Call the repository method to get user public key
	user, err := userGrpcServer.userRepo.GetByEmail(ctx, req.Msg.Email)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get user public key: %w", err))
	}

	resp := &userv1.GetUserPublicKeyResponse{
		UserId:    user.ID.Hex(),
		PublicKey: *user.KeySet.PublicKey,
	}

	return connect.NewResponse(resp), nil
}
