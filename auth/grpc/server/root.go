package server

import "github.com/atomic-blend/backend/shared/repositories/user"

// UserGrpcServer is the gRPC server for user-related operations
type UserGrpcServer struct {
	userRepo 	user.UserRepositoryInterface
}

// NewUserGrpcServer creates a new UserGrpcServer instance
func NewUserGrpcServer(userRepo user.UserRepositoryInterface) *UserGrpcServer {
	return &UserGrpcServer{
		userRepo: userRepo,
	}
}