package service

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/gnad103/go-ex/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServiceServer implements the UserService gRPC service
type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	mu    sync.Mutex
	users map[string]*pb.User
}

// NewUserServiceServer creates a new UserServiceServer
func NewUserServiceServer() *UserServiceServer {
	return &UserServiceServer{
		users: make(map[string]*pb.User),
	}
}

// GetUser retrieves a user by ID
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user with ID %s not found", req.Id)
	}

	return user, nil
}

// CreateUser creates a new user
func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a new user with a simple ID generation
	id := fmt.Sprintf("user-%d", len(s.users)+1)
	user := &pb.User{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	// Store the user
	s.users[id] = user

	return user, nil
}

// ListUsers retrieves all users
func (s *UserServiceServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	response := &pb.ListUsersResponse{}
	for _, user := range s.users {
		response.Users = append(response.Users, user)
	}

	return response, nil
}
