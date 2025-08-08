package grpserver

import "github.com/atomic-blend/backend/mail/repositories"

// GrpcServer is the gRPC server for the productivity service
type GrpcServer struct {
	mailRepo     repositories.MailRepositoryInterface
	sendMailRepo repositories.SendMailRepositoryInterface
}

// NewGrpcServer create a new instance of GrpcServer
func NewGrpcServer(mailRepo repositories.MailRepositoryInterface, sendMailRepo repositories.SendMailRepositoryInterface) *GrpcServer {
	return &GrpcServer{
		mailRepo: mailRepo,
	}
}
