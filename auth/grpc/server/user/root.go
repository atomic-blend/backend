package server

import "github.com/atomic-blend/backend/auth/repositories"

// UserGrpcServer is the gRPC server for user-related operations
type UserGrpcServer struct {
	userRepo 	repositories.UserRepositoryInterface
}

// NewUserGrpcServer creates a new UserGrpcServer instance
func NewUserGrpcServer(userRepo repositories.UserRepositoryInterface) *UserGrpcServer {
	return &UserGrpcServer{
		userRepo: userRepo,
	}
}