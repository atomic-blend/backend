// Package server provides the gRPC server for the mail-server service
package server

// GrpcServer is the gRPC server for the mail-server service
type GrpcServer struct {
}

// NewGrpcServer create a new instance of GrpcServer
func NewGrpcServer() *GrpcServer {
	return &GrpcServer{}
}
