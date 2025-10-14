package server

// GrpcServer is the gRPC server for the productivity service
type GrpcServer struct {
}

// NewGrpcServer create a new instance of GrpcServer
func NewGrpcServer() *GrpcServer {
	return &GrpcServer{}
}
