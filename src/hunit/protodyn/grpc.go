// Routines for interacting with gRPC services
//
//	func main() {
//		// Load service definitions
//		registry := NewServiceRegistry()
//
//		// Load from FileDescriptorSet
//		fds := &descriptorpb.FileDescriptorSet{} // loaded from file
//		err := registry.LoadFileDescriptorSet(fds)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		// Create gRPC connection
//		conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
//		if err != nil {
//			log.Fatal(err)
//		}
//		defer conn.Close()
//
//		// Create client
//		client := NewClient(conn, registry)
//
//		// Make a call
//		requestJSON := []byte(`{"name": "John", "age": 30}`)
//		responseMessage, err := client.Call(
//			context.Background(),
//			"com.example.UserService",
//			"GetUser",
//			requestJSON,
//			&CallOptions{},
//		)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		responseJSON, err := MarshalJSON(responseMessage)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		fmt.Println("Response:", string(responseJSON))
//	}
package protodyn

import (
	"context"
	"fmt"
	"io/ioutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

// ServiceRegistry manages protobuf service descriptors
type ServiceRegistry struct {
	files    *protoregistry.Files
	services map[string]protoreflect.ServiceDescriptor
}

// NewServiceRegistry creates a new service registry
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		files:    &protoregistry.Files{},
		services: make(map[string]protoreflect.ServiceDescriptor),
	}
}

// LoadFileDescriptorSetFromPath loads service definitions from a FileDescriptorSet found at the specified path
func (r *ServiceRegistry) LoadFileDescriptorSetFromPath(p string) error {
	descBytes, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	var fds descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(descBytes, &fds); err != nil {
		return err
	}

	return r.LoadFileDescriptorSet(&fds)
}

// LoadFileDescriptorSet loads service definitions from a FileDescriptorSet
func (r *ServiceRegistry) LoadFileDescriptorSet(fds *descriptorpb.FileDescriptorSet) error {
	// Use protodesc.NewFiles to create a complete registry from the FileDescriptorSet
	files, err := protodesc.NewFiles(fds)
	if err != nil {
		return fmt.Errorf("failed to create file registry: %w", err)
	}

	// Copy all files to our registry
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		if err := r.files.RegisterFile(fd); err != nil {
			// Skip files that are already registered (e.g., well-known types)
			return true
		}
		return true
	})

	// Collect all services from the registered files
	r.files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		services := fd.Services()
		for i := 0; i < services.Len(); i++ {
			svc := services.Get(i)
			r.services[string(svc.FullName())] = svc
		}
		return true
	})

	return nil
}

// GetService retrieves a service descriptor by full name
func (r *ServiceRegistry) GetService(fullName string) (protoreflect.ServiceDescriptor, error) {
	svc, exists := r.services[fullName]
	if !exists {
		return nil, fmt.Errorf("service %s not found", fullName)
	}
	return svc, nil
}

// ListServices returns all registered service names
func (r *ServiceRegistry) ListServices() []string {
	names := make([]string, 0, len(r.services))
	for name := range r.services {
		names = append(names, name)
	}
	return names
}

// Client represents a dynamic gRPC client
type Client struct {
	conn     *grpc.ClientConn
	registry *ServiceRegistry
}

// NewClient creates a new dynamic gRPC client
func NewClient(conn *grpc.ClientConn, registry *ServiceRegistry) *Client {
	return &Client{
		conn:     conn,
		registry: registry,
	}
}

// CallOptions contains options for gRPC calls
type CallOptions struct {
	// Additional gRPC call options can be added here
	GrpcOptions []grpc.CallOption
}

func (c *Client) Method(ctx context.Context, serviceName, methodName string) (protoreflect.MethodDescriptor, error) {
	svc, err := c.registry.GetService(serviceName)
	if err != nil {
		return nil, fmt.Errorf("Service is not registered: %w", err)
	}

	method := svc.Methods().ByName(protoreflect.Name(methodName))
	if method == nil {
		return nil, fmt.Errorf("Method %s not found in service %s", methodName, serviceName)
	}

	return method, nil
}

// Call performs a unary RPC call
func (c *Client) Call(cxt context.Context, serviceName, methodName string, requestJSON []byte, opts *CallOptions) (*dynamicpb.Message, error) {
	method, err := c.Method(cxt, serviceName, methodName)
	if err != nil {
		return nil, err
	}

	// Create request message
	reqmsg := dynamicpb.NewMessage(method.Input())
	// unmarshal request payload from JSON
	err = UnmarshalJSON(requestJSON, reqmsg)
	if err != nil {
		return nil, fmt.Errorf("could not encode request: %w", err)
	}

	// Create response message
	rspmsg := dynamicpb.NewMessage(method.Output())
	// Build method path
	methodPath := fmt.Sprintf("/%s/%s", serviceName, methodName)

	// Prepare gRPC options
	grpcOpts := []grpc.CallOption{}
	if opts != nil && opts.GrpcOptions != nil {
		grpcOpts = append(grpcOpts, opts.GrpcOptions...)
	}

	// Make the call
	err = c.conn.Invoke(cxt, methodPath, reqmsg, rspmsg, grpcOpts...)
	if err != nil {
		return nil, c.formatGrpcError(err)
	}

	return rspmsg, nil
}

// formatGrpcError formats gRPC errors with human-readable descriptions
func (c *Client) formatGrpcError(err error) error {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.NotFound:
			return fmt.Errorf("service or method not found: %s", st.Message())
		case codes.InvalidArgument:
			return fmt.Errorf("invalid request: %s", st.Message())
		case codes.Unauthenticated:
			return fmt.Errorf("authentication required: %s", st.Message())
		case codes.PermissionDenied:
			return fmt.Errorf("permission denied: %s", st.Message())
		case codes.Unavailable:
			return fmt.Errorf("service unavailable: %s", st.Message())
		case codes.DeadlineExceeded:
			return fmt.Errorf("request timeout: %s", st.Message())
		default:
			return fmt.Errorf("gRPC error (%s): %s", st.Code().String(), st.Message())
		}
	}
	return err
}
