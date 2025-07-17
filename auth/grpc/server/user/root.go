package server

import "github.com/atomic-blend/backend/auth/repositories"

type UserGrpcServer struct {
	userRepo 	repositories.UserRepositoryInterface
}

func NewUserGrpcServer(userRepo repositories.UserRepositoryInterface) *UserGrpcServer {
	return &UserGrpcServer{
		userRepo: userRepo,
	}
}