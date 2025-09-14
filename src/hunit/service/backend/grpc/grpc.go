package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/instaunit/instaunit/hunit/expr"
	"github.com/instaunit/instaunit/hunit/expr/runtime"
	"github.com/instaunit/instaunit/hunit/protodyn"
	"github.com/instaunit/instaunit/hunit/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/dynamicpb"
)

var marshalOptions = protojson.MarshalOptions{}

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
	svcreg *protodyn.ServiceRegistry
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

	var (
		reg  = protodyn.NewServiceRegistry()
		root = path.Dir(conf.Path)
	)
	for _, p := range suite.Protos {
		err := reg.LoadFileDescriptorSetFromPath(path.Join(root, p))
		if err != nil {
			return nil, fmt.Errorf("Could not load Protobuf descriptor set: %v: %w", p, err)
		}
	}

	svc := &grpcService{
		conf:   conf,
		suite:  suite,
		svcreg: reg,
		vars: expr.Variables{
			"std": runtime.Stdlib,
		},
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

	method, err := s.svcreg.MethodForName(mname)
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("Method is not defined in any known service: %v", mname))
	}

	// Create request message to receive the incoming data
	reqmsg := dynamicpb.NewMessage(method.Input())
	// Receive the request (this reads the request but we don't need to use it for static responses)
	err = stream.RecvMsg(reqmsg)
	if err != nil {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("Could not receive request message: %v", err))
	}

	// Convert the request to JSON data...
	reqdata, err := marshalOptions.Marshal(reqmsg)
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("Could not marshal request message: %v", err))
	}
	// ...and unmarshal back to generic types for us in interpolation vars
	var reqjson interface{}
	err = json.Unmarshal(reqdata, &reqjson)
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("Could not convert request to JSON: %v", err))
	}

	epoint, ok := s.suite.MatchEndpoint(mname, reqmsg)
	if !ok {
		return status.Error(codes.Unimplemented, fmt.Sprintf("Not implemented: %s", mname))
	}
	rsp := epoint.Response
	if rsp == nil {
		return fmt.Errorf("Endpoint has no response defined")
	}

	// interpolate the response
	vars := s.vars.With(expr.Variables{
		"method":   mname,
		"endpoint": epoint,
		"request": expr.Variables{
			"value": reqjson,
		},
	})
	spew.Dump(vars)
	rspdata, err := expr.Interpolate(rsp.Entity, vars)
	if err != nil {
		return fmt.Errorf("Could not interpolate response: %w", err)
	}

	// Create request message to receive the incoming data
	rspmsg := dynamicpb.NewMessage(method.Output())
	// Use protojson to unmarshal into the dynamic message
	err = protojson.Unmarshal([]byte(rspdata), rspmsg)
	if err != nil {
		return fmt.Errorf("Could not unmarshal JSON response to protobuf: %v", err)
	}
	// Send our response
	err = stream.SendMsg(rspmsg)
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("Could not send response message: %v", err))
	}

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
