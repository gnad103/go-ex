package service

import (
	"context"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gnad103/go-ex/proto"
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
