package security

type GrpcClient struct {
}

type GrpcServer struct {
}

func NewGrpcClient() *GrpcClient {
	return &GrpcClient{}
}

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{}
}
