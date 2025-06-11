package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/instaunit/grpcexample/service/proto"

	"google.golang.org/grpc"
)

type userServer struct {
	proto.UnimplementedUserServiceServer
	users  map[int32]*proto.User
	nextID int32
}

func newUserServer() *userServer {
	return &userServer{
		users:  make(map[int32]*proto.User),
		nextID: 1,
	}
}

func (s *userServer) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	user, exists := s.users[req.UserId]
	if !exists {
		return nil, fmt.Errorf("user %d not found", req.UserId)
	}

	return &proto.GetUserResponse{User: user}, nil
}

func (s *userServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	user := req.User
	user.Id = s.nextID
	s.nextID++

	s.users[user.Id] = user

	return &proto.CreateUserResponse{User: user}, nil
}

func (s *userServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	_, exists := s.users[req.UserId]
	if !exists {
		return nil, fmt.Errorf("user %d not found", req.UserId)
	}

	user := req.User
	user.Id = req.UserId
	s.users[req.UserId] = user

	return &proto.UpdateUserResponse{User: user}, nil
}

func (s *userServer) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	_, exists := s.users[req.UserId]
	if !exists {
		return &proto.DeleteUserResponse{Success: false}, nil
	}

	delete(s.users, req.UserId)
	return &proto.DeleteUserResponse{Success: true}, nil
}

func (s *userServer) ListUsers(ctx context.Context, req *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
	var users []*proto.User
	for _, user := range s.users {
		// Simple type filtering
		if len(req.Types) > 0 {
			found := false
			for _, t := range req.Types {
				if user.Type == t {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		users = append(users, user)
	}

	return &proto.ListUsersResponse{
		Users:      users,
		TotalCount: int32(len(users)),
	}, nil
}

type testServer struct {
	proto.UnimplementedTestServiceServer
}

func (s *testServer) TestNumericTypes(ctx context.Context, req *proto.NumericTestRequest) (*proto.NumericTestResponse, error) {
	// Echo back the data with a message
	return &proto.NumericTestResponse{
		Data:    req.Data,
		Message: "All numeric types processed successfully",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8800")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// Register services
	proto.RegisterUserServiceServer(s, newUserServer())
	proto.RegisterTestServiceServer(s, &testServer{})

	log.Println("gRPC server listening on port 8800")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
