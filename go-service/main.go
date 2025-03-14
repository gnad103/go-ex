package main

import (
	"log"
	"net"

	"github.com/gnad103/go-ex/go-service/service"
	pb "github.com/gnad103/go-ex/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Create a TCP listener on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Create and register our user service implementation
	userService := service.NewUserServiceServer()
	pb.RegisterUserServiceServer(grpcServer, userService)

	// Register reflection service on gRPC server to allow for service discovery
	reflection.Register(grpcServer)

	log.Println("Go User Service is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
