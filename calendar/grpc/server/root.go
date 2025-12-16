package global

import (
	userclient "github.com/atomic-blend/backend/shared/grpc/user"
)

// GrpcServer is the gRPC server for the calendar service
type GrpcServer struct {
	UserClient userclient.Interface
}

// NewGrpcServer create a new instance of GrpcServer
func NewGrpcServer() *GrpcServer {
	return &GrpcServer{}
}
