package server

import "github.com/atomic-blend/backend/shared/repositories/user"

// UserGrpcServer is the gRPC server for user-related operations
type UserGrpcServer struct {
	userRepo user.Interface
}

// NewUserGrpcServer creates a new UserGrpcServer instance
func NewUserGrpcServer(userRepo user.Interface) *UserGrpcServer {
	return &UserGrpcServer{
		userRepo: userRepo,
	}
}
