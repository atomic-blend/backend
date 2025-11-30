package global

// GrpcServer is the gRPC server for the calendar service
type GrpcServer struct {
}

// NewGrpcServer create a new instance of GrpcServer
func NewGrpcServer() *GrpcServer {
	return &GrpcServer{}
}
