package user

import (
	"github.com/atomic-blend/backend/productivity/repositories"
)

// grpcServer is the gRPC server for the productivity service
type grpcServer struct {
	taskRepo repositories.TaskRepositoryInterface
}

// NewGrpcServer create a new instance of GrpcServer
func NewGrpcServer(taskRepo repositories.TaskRepositoryInterface) *grpcServer {
	return &grpcServer{
		taskRepo: taskRepo,
	}
}
