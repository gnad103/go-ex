package main

import (
	"log"
	"net"

	"github.com/yourusername/go-ex/proto"
	"github.com/yourusername/go-ex/user-service/service"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterUserServiceServer(s, &service.UserServer{})

	log.Println("User service is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
