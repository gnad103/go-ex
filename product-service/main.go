package main

import (
	"log"
	"net"

	"github.com/yourusername/go-ex/product-service/service"
	"github.com/yourusername/go-ex/proto"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a connection to the user service
	userConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	defer userConn.Close()

	userClient := proto.NewUserServiceClient(userConn)

	s := grpc.NewServer()
	proto.RegisterProductServiceServer(s, service.NewProductServer(userClient))

	log.Println("Product service is running on port 50052...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
