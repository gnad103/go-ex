package service

import (
	"context"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yourusername/go-ex/proto"
)

type ProductServer struct {
	proto.UnimplementedProductServiceServer
	mu         sync.Mutex
	products   map[int64]*proto.ProductResponse
	nextID     int64
	userClient proto.UserServiceClient
}

func NewProductServer(userClient proto.UserServiceClient) *ProductServer {
	return &ProductServer{
		products:   make(map[int64]*proto.ProductResponse),
		nextID:     1,
		userClient: userClient,
	}
}

func (s *ProductServer) GetProduct(ctx context.Context, req *proto.ProductRequest) (*proto.ProductResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	product, exists := s.products[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "product with ID %d not found", req.Id)
	}

	return product, nil
}

func (s *ProductServer) CreateProduct(ctx context.Context, req *proto.CreateProductRequest) (*proto.ProductResponse, error) {
	// First verify that the user exists
	_, err := s.userClient.GetUser(ctx, &proto.UserRequest{Id: req.UserId})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "user with ID %d not found: %v", req.UserId, err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	product := &proto.ProductResponse{
		Id:          s.nextID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		UserId:      req.UserId,
	}

	s.products[s.nextID] = product
	s.nextID++

	return product, nil
}

func (s *ProductServer) GetProductsForUser(ctx context.Context, req *proto.UserProductRequest) (*proto.ProductListResponse, error) {
	// First verify that the user exists
	_, err := s.userClient.GetUser(ctx, &proto.UserRequest{Id: req.UserId})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "user with ID %d not found: %v", req.UserId, err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var userProducts []*proto.ProductResponse

	for _, product := range s.products {
		if product.UserId == req.UserId {
			userProducts = append(userProducts, product)
		}
	}

	return &proto.ProductListResponse{Products: userProducts}, nil
}
package main

import (
	"context"
	"log"
	"time"

	"github.com/yourusername/go-ex/proto"
	"google.golang.org/grpc"
)

func main() {
	// Connect to user service
	userConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	defer userConn.Close()
	
	userClient := proto.NewUserServiceClient(userConn)
	
	// Connect to product service
	productConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to product service: %v", err)
	}
	defer productConn.Close()
	
	productClient := proto.NewProductServiceClient(productConn)
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
	// Create a user
	user, err := userClient.CreateUser(ctx, &proto.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	})
	if err != nil {
		log.Fatalf("Could not create user: %v", err)
	}
	log.Printf("Created user: %v", user)
	
	// Create a product for the user
	product, err := productClient.CreateProduct(ctx, &proto.CreateProductRequest{
		Name:        "Laptop",
		Description: "High-performance laptop",
		Price:       999.99,
		UserId:      user.Id,
	})
	if err != nil {
		log.Fatalf("Could not create product: %v", err)
	}
	log.Printf("Created product: %v", product)
	
	// Get products for the user
	products, err := productClient.GetProductsForUser(ctx, &proto.UserProductRequest{
		UserId: user.Id,
	})
	if err != nil {
		log.Fatalf("Could not get products for user: %v", err)
	}
	log.Printf("Products for user %d: %v", user.Id, products)
}