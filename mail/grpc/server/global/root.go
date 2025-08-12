package global

import "github.com/atomic-blend/backend/mail/repositories"

// GrpcServer is the gRPC server for the productivity service
type GrpcServer struct {
	sendMailRepository repositories.SendMailRepositoryInterface
}

// NewGrpcServer create a new instance of GrpcServer
func NewGrpcServer(sendMailRepository repositories.SendMailRepositoryInterface) *GrpcServer {
	return &GrpcServer{
		sendMailRepository: sendMailRepository,
	}
}
