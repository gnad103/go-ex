package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/gnad103/go-ex/proto"
	"github.com/gnad103/go-ex/user-service/service"
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
