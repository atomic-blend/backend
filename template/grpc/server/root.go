package global

// GrpcServer is the gRPC server for the template service
type GrpcServer struct {
}

// NewGrpcServer create a new instance of GrpcServer
func NewGrpcServer() *GrpcServer {
	return &GrpcServer{}
}
