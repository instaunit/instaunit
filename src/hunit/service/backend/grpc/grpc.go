package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/instaunit/instaunit/hunit/expr"
	"github.com/instaunit/instaunit/hunit/expr/runtime"
	"github.com/instaunit/instaunit/hunit/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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
	var (
		sinfo = grpc.ServerTransportStreamFromContext(stream.Context())
		mname = strings.TrimPrefix(sinfo.Method(), "/")
	)

	epoint, ok := s.suite.FindEndpoint(mname)
	if !ok {
		return status.Error(codes.Unimplemented, fmt.Sprintf("Not implemented: %s", mname))
	}

	log.Printf("Handling call to method: %s: %v", mname, epoint)

	// log.Printf("Successfully handled call to %s", fullMethodName)
	return nil
}

func (s *grpcService) Start() error {
	lnr, err := net.Listen("tcp", s.conf.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", s.conf.Addr, err)
	}

	go func() {
		err = s.server.Serve(lnr)
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
