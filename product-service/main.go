package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/gnad103/go-ex/product-service/service"
	"github.com/gnad103/go-ex/proto"
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
