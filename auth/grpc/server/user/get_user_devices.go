package server

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/atomic-blend/backend/grpc/gen/auth"
)

// GetUserDevices is the gRPC method to retrieve user devices
func (userGrpcServer *UserGrpcServer) GetUserDevices(ctx context.Context, req *connect.Request[auth.GetUserDevicesRequest]) (*connect.Response[auth.GetUserDevicesResponse], error) {
	// Call the repository method to get user devices
	user, err := userGrpcServer.userRepo.GetByID(ctx, req.Msg.User.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get user devices: %w", err))
	}

	devices := []*auth.UserDevice{}
	for _, device := range user.Devices {
		deviceTz := ""
		if device.DeviceTimezone != nil {
			deviceTz = *device.DeviceTimezone
		}
		devices = append(devices, &auth.UserDevice{
			DeviceID:          device.DeviceID,
			DeviceName:        device.DeviceName,
			FcmToken:         device.FcmToken,
			DeviceTimezone:      deviceTz,   
		})
	}

	// Create the response
	resp := &auth.GetUserDevicesResponse{
		Devices: devices,
	}

	return connect.NewResponse(resp), nil
}
