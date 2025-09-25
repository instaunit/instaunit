package grpc

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"
	"strings"

	"github.com/instaunit/instaunit/hunit/expr"
	"github.com/instaunit/instaunit/hunit/expr/runtime"
	"github.com/instaunit/instaunit/hunit/protodyn"
	"github.com/instaunit/instaunit/hunit/service"
	"github.com/instaunit/instaunit/hunit/service/status"
	"github.com/instaunit/instaunit/hunit/text"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
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

	conf.Status.Set(conf.Addr, status.Pending)

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
		grpc.UnknownServiceHandler(svc.handleUnknownService),
	)

	// Enable server reflection for debugging
	reflection.Register(grs)
	// update the service
	svc.server = grs

	return svc, nil
}

// handleUnknownService handles all incoming gRPC calls
func (s *grpcService) handleUnknownService(srv interface{}, stream grpc.ServerStream) error {
	var (
		sinfo = grpc.ServerTransportStreamFromContext(stream.Context())
		mname = strings.TrimPrefix(sinfo.Method(), "/")
	)
	method, err := s.svcreg.MethodForName(mname)
	if err != nil {
		return grpcErrf(codes.Internal, "Method is not defined in any known service: %v", mname)
	}

	// Create request message to receive the incoming data
	msgdsc := method.Input()
	reqmsg := dynamicpb.NewMessage(msgdsc)
	// Receive the request (this reads the request but we don't need to use it for static responses)
	err = stream.RecvMsg(reqmsg)
	if err != nil {
		return grpcErrf(codes.InvalidArgument, "Could not receive request message: %v: %v", msgdsc.FullName(), err)
	}

	// Convert the request to JSON data...
	reqdata, err := marshalOptions.Marshal(reqmsg)
	if err != nil {
		return grpcErrf(codes.Internal, "Could not marshal request message: %v", err)
	}
	// ...and unmarshal back to generic types for us in interpolation vars
	var reqjson interface{}
	err = json.Unmarshal(reqdata, &reqjson)
	if err != nil {
		return grpcErrf(codes.Internal, "Could not convert request to JSON: %v", err)
	}

	epoint, ok := s.suite.MatchEndpoint(mname, reqmsg)
	if !ok {
		return grpcErrf(codes.Unimplemented, "Not implemented: %s", mname)
	}
	rsp := epoint.Response
	if rsp == nil {
		return grpcErrf(codes.Internal, "Endpoint has no response defined")
	}

	// interpolate the response
	vars := s.vars.With(expr.Variables{
		"method":   mname,
		"endpoint": epoint,
		"request": expr.Variables{
			"value": reqjson,
		},
	})

	if e := rsp.Error; e != nil && e.Code != "" {
		return grpcErrf(parseErrorCode(e.Code, codes.Unknown), text.Coalesce(e.Message, e.Code))
	} else if rsp.Status != 0 {
		return grpcErrf(codes.Code(rsp.Status), "Status")
	}

	rspdata, err := expr.Interpolate(rsp.Entity, vars)
	if err != nil {
		return grpcErrf(codes.Internal, "Could not interpolate response: %v", err)
	}

	// Create request message to receive the incoming data
	msgdsc = method.Output()
	rspmsg := dynamicpb.NewMessage(msgdsc)
	// Use protojson to unmarshal into the dynamic message
	err = protojson.Unmarshal([]byte(rspdata), rspmsg)
	if err != nil {
		return grpcErrf(codes.Internal, "Could not unmarshal JSON response to protobuf: %v: %v (this is a configuration error in your Instaunit service; make sure your response conforms to the endpoint's output type)", msgdsc.FullName(), err)
	}

	// Send our response
	err = stream.SendMsg(rspmsg)
	if err != nil {
		return grpcErrf(codes.Internal, "Could not send response message: %v", err)
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

	s.conf.Status.Set(s.conf.Addr, status.Ready)
	return nil
}

func (s *grpcService) Stop() error {
	s.conf.Status.Set(s.conf.Addr, status.Stopped)
	s.server.GracefulStop()
	return nil
}
