package user

// grpcServer is the gRPC server for the productivity service
type grpcServer struct{}

// NewGrpcServer create a new instance of GrpcServer
func NewGrpcServer() *grpcServer {
	return &grpcServer{}
}
