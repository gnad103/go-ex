package service

import (
	"context"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gnad103/go-ex/proto"
)

type UserServer struct {
	proto.UnimplementedUserServiceServer
	mu     sync.Mutex
	users  map[int64]*proto.UserResponse
	nextID int64
}

func NewUserServer() *UserServer {
	return &UserServer{
		users:  make(map[int64]*proto.UserResponse),
		nextID: 1,
	}
}

func (s *UserServer) GetUser(ctx context.Context, req *proto.UserRequest) (*proto.UserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user with ID %d not found", req.Id)
	}

	return user, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.UserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := &proto.UserResponse{
		Id:    s.nextID,
		Name:  req.Name,
		Email: req.Email,
	}

	s.users[s.nextID] = user
	s.nextID++

	return user, nil
}
