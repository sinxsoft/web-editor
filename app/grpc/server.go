package grpc

// server.go

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	//pb "helloworld/helloworld"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterGreeterServer(s, &server{})
	s.Serve(lis)
}
