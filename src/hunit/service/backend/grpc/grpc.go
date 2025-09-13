package grpc

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/instaunit/instaunit/hunit/expr"
	"github.com/instaunit/instaunit/hunit/expr/runtime"
	"github.com/instaunit/instaunit/hunit/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func logln(v ...any) {
	fmt.Fprintln(os.Stderr, v...)
}

func logf(f string, a ...any) {
	if n := len(f); n == 0 {
		fmt.Fprintln(os.Stderr)
	} else if f[n] == '\n' {
		fmt.Fprintf(os.Stderr, f, a...)
	} else {
		fmt.Fprintf(os.Stderr, f, a...)
		fmt.Fprintln(os.Stderr)
	}
}

// gRPC service
type grpcService struct {
	conf   service.Config
	suite  *Suite
	server *grpc.Server
	vars   expr.Variables
}

// Create a new service
func New(conf service.Config) (service.Service, error) {
	src, err := os.Open(conf.Path)
	if err != nil {
		return nil, err
	}
	defer src.Close()
	suite, err := LoadSuite(src)
	if err != nil {
		return nil, err
	}

	vars := expr.Variables{
		"std": runtime.Stdlib,
	}

	svc := &grpcService{
		conf:  conf,
		suite: suite,
		vars:  vars,
	}

	// Create gRPC server with unknown service handler
	grs := grpc.NewServer(
		grpc.UnaryInterceptor(svc.unaryInterceptor),
		grpc.UnknownServiceHandler(svc.handleUnknownService),
	)

	// Enable server reflection for debugging
	reflection.Register(grs)
	// update the service
	svc.server = grs

	return svc, nil
}

func (s *grpcService) unaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	return nil, nil
}

// handleUnknownService handles all incoming gRPC calls
func (s *grpcService) handleUnknownService(srv interface{}, stream grpc.ServerStream) error {
	fmt.Println(">>>>>>>>>>>>>>>>> ZZ", srv)

	// Get the method name from the context
	// fullMethodName, ok := grpc.MethodFromContext(stream.Context())
	// if !ok {
	// 	return status.Error(codes.Internal, "failed to get method name from context")
	// }

	// // Remove leading slash from method name (e.g., "/package.Service/Method" -> "package.Service/Method")
	// if strings.HasPrefix(fullMethodName, "/") {
	// 	fullMethodName = fullMethodName[1:]
	// }

	// log.Printf("Handling call to method: %s", fullMethodName)

	// log.Printf("Successfully handled call to %s", fullMethodName)
	return nil
}

func (s *grpcService) Start() error {
	port := 31222

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %v", port, err)
	}

	go func() {
		err = s.server.Serve(listener)
		if err != nil {
			logf("gRPC server error: %v", err)
		}
	}()

	return nil
}

func (s *grpcService) Stop() error {
	s.server.GracefulStop()
	return nil
}

// populateMessageFromConfig populates a protobuf message from configuration data
// func (m *MockGRPCServer) populateMessageFromConfig(msg *dynamicpb.Message, configData map[string]interface{}) error {
// 	// Convert config data to JSON for easier protobuf unmarshaling
// 	jsonData, err := json.Marshal(configData)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal config to JSON: %v", err)
// 	}
//
// 	// Use protojson to unmarshal into the dynamic message
// 	if err := protojson.Unmarshal(jsonData, msg); err != nil {
// 		return fmt.Errorf("failed to unmarshal JSON to protobuf: %v", err)
// 	}
//
// 	return nil
// }
