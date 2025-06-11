package dynrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

// ValidationError represents a collection of validation errors
type ValidationError struct {
	Errors []FieldError
}

func (e *ValidationError) Error() string {
	if len(e.Errors) == 0 {
		return "validation failed"
	}
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}

	var sb strings.Builder
	sb.WriteString("validation failed with multiple errors:\n")
	for _, err := range e.Errors {
		sb.WriteString("  - ")
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}
	return sb.String()
}

// FieldError represents a validation error for a specific field
type FieldError struct {
	Field    string // dot-separated field path
	Expected string // expected type/format
	Actual   string // actual value or type received
	Message  string // human-readable description
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("field '%s': %s (expected %s, got %s)",
		e.Field, e.Message, e.Expected, e.Actual)
}

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
	if err := c.jsonToProto(requestJSON, reqmsg, ""); err != nil {
		return nil, fmt.Errorf("request encoding failed: %w", err)
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

// jsonToProto converts JSON to a protobuf message with full validation
func (c *Client) jsonToProto(jsonData []byte, msg *dynamicpb.Message, fieldPath string) error {
	var jsonObj map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonObj); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	var validationErrors []FieldError
	c.populateMessage(msg, jsonObj, fieldPath, &validationErrors)

	if len(validationErrors) > 0 {
		return &ValidationError{Errors: validationErrors}
	}

	return nil
}

// populateMessage populates a protobuf message from a JSON object
func (c *Client) populateMessage(msg *dynamicpb.Message, jsonObj map[string]interface{}, basePath string, errors *[]FieldError) {
	desc := msg.Descriptor()
	fields := desc.Fields()

	// Track oneof fields to validate exclusivity
	oneofFields := make(map[protoreflect.OneofDescriptor][]string)

	// Process each JSON field
	for jsonKey, jsonValue := range jsonObj {
		fieldPath := c.buildFieldPath(basePath, jsonKey)

		// Find the protobuf field
		field := fields.ByJSONName(jsonKey)
		if field == nil {
			field = fields.ByName(protoreflect.Name(jsonKey))
		}

		if field == nil {
			*errors = append(*errors, FieldError{
				Field:    fieldPath,
				Expected: "known field",
				Actual:   jsonKey,
				Message:  "unknown field",
			})
			continue
		}

		// Handle oneof fields
		if field.ContainingOneof() != nil {
			oneofFields[field.ContainingOneof()] = append(
				oneofFields[field.ContainingOneof()], jsonKey)
		}

		// Convert and set the field value
		c.setFieldValue(msg, field, jsonValue, fieldPath, errors)
	}

	// Validate oneof constraints
	for oneof, setFields := range oneofFields {
		if len(setFields) > 1 {
			for _, fieldName := range setFields {
				fieldPath := c.buildFieldPath(basePath, fieldName)
				*errors = append(*errors, FieldError{
					Field:    fieldPath,
					Expected: "only one field in oneof group",
					Actual:   fmt.Sprintf("multiple fields set: %v", setFields),
					Message:  fmt.Sprintf("oneof group '%s' has multiple fields set", oneof.Name()),
				})
			}
		}
	}
}

// setFieldValue sets a field value with type validation
func (c *Client) setFieldValue(msg *dynamicpb.Message, field protoreflect.FieldDescriptor, value interface{}, fieldPath string, errors *[]FieldError) {
	if field.IsList() {
		c.setRepeatedField(msg, field, value, fieldPath, errors)
	} else if field.IsMap() {
		c.setMapField(msg, field, value, fieldPath, errors)
	} else {
		c.setSingleField(msg, field, value, fieldPath, errors)
	}
}

// setSingleField sets a single (non-repeated, non-map) field value
func (c *Client) setSingleField(msg *dynamicpb.Message, field protoreflect.FieldDescriptor,
	value interface{}, fieldPath string, errors *[]FieldError) {

	convertedValue, err := c.convertValue(value, field, fieldPath)
	if err != nil {
		*errors = append(*errors, *err)
		return
	}

	msg.Set(field, convertedValue)
}

// setRepeatedField sets a repeated field value
func (c *Client) setRepeatedField(msg *dynamicpb.Message, field protoreflect.FieldDescriptor,
	value interface{}, fieldPath string, errors *[]FieldError) {

	// Handle null or empty arrays
	if value == nil {
		msg.Set(field, msg.NewField(field))
		return
	}

	jsonArray, ok := value.([]interface{})
	if !ok {
		*errors = append(*errors, FieldError{
			Field:    fieldPath,
			Expected: "array",
			Actual:   fmt.Sprintf("%T", value),
			Message:  "repeated field must be an array",
		})
		return
	}

	list := msg.Mutable(field).List()
	for i, item := range jsonArray {
		itemPath := fmt.Sprintf("%s[%d]", fieldPath, i)
		convertedItem, err := c.convertValue(item, field, itemPath)
		if err != nil {
			*errors = append(*errors, *err)
			continue
		}
		list.Append(convertedItem)
	}
}

// setMapField sets a map field value
func (c *Client) setMapField(msg *dynamicpb.Message, field protoreflect.FieldDescriptor,
	value interface{}, fieldPath string, errors *[]FieldError) {

	if value == nil {
		msg.Set(field, msg.NewField(field))
		return
	}

	jsonMap, ok := value.(map[string]interface{})
	if !ok {
		*errors = append(*errors, FieldError{
			Field:    fieldPath,
			Expected: "object",
			Actual:   fmt.Sprintf("%T", value),
			Message:  "map field must be an object",
		})
		return
	}

	mapValue := msg.Mutable(field).Map()
	keyField := field.MapKey()
	valueField := field.MapValue()

	for k, v := range jsonMap {
		keyPath := fmt.Sprintf("%s.%s", fieldPath, k)

		// Convert key
		convertedKey, err := c.convertValue(k, keyField, keyPath)
		if err != nil {
			*errors = append(*errors, *err)
			continue
		}

		// Convert value
		convertedValue, err := c.convertValue(v, valueField, keyPath)
		if err != nil {
			*errors = append(*errors, *err)
			continue
		}

		mapValue.Set(convertedKey.MapKey(), convertedValue)
	}
}

// convertValue converts a JSON value to the appropriate protobuf type
func (c *Client) convertValue(value interface{}, field protoreflect.FieldDescriptor,
	fieldPath string) (protoreflect.Value, *FieldError) {

	switch field.Kind() {
	case protoreflect.BoolKind:
		if b, ok := value.(bool); ok {
			return protoreflect.ValueOfBool(b), nil
		}
		return protoreflect.Value{}, &FieldError{
			Field:    fieldPath,
			Expected: "boolean",
			Actual:   fmt.Sprintf("%T", value),
			Message:  "expected boolean value",
		}

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		if num, ok := value.(float64); ok {
			return protoreflect.ValueOfInt32(int32(num)), nil
		}
		return protoreflect.Value{}, &FieldError{
			Field:    fieldPath,
			Expected: "int32",
			Actual:   fmt.Sprintf("%T", value),
			Message:  "expected 32-bit integer",
		}

	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		if num, ok := value.(float64); ok {
			return protoreflect.ValueOfInt64(int64(num)), nil
		}
		return protoreflect.Value{}, &FieldError{
			Field:    fieldPath,
			Expected: "int64",
			Actual:   fmt.Sprintf("%T", value),
			Message:  "expected 64-bit integer",
		}

	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		if num, ok := value.(float64); ok {
			return protoreflect.ValueOfUint32(uint32(num)), nil
		}
		return protoreflect.Value{}, &FieldError{
			Field:    fieldPath,
			Expected: "uint32",
			Actual:   fmt.Sprintf("%T", value),
			Message:  "expected 32-bit unsigned integer",
		}

	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		if num, ok := value.(float64); ok {
			return protoreflect.ValueOfUint64(uint64(num)), nil
		}
		return protoreflect.Value{}, &FieldError{
			Field:    fieldPath,
			Expected: "uint64",
			Actual:   fmt.Sprintf("%T", value),
			Message:  "expected 64-bit unsigned integer",
		}

	case protoreflect.FloatKind:
		if num, ok := value.(float64); ok {
			return protoreflect.ValueOfFloat32(float32(num)), nil
		}
		return protoreflect.Value{}, &FieldError{
			Field:    fieldPath,
			Expected: "float32",
			Actual:   fmt.Sprintf("%T", value),
			Message:  "expected 32-bit float",
		}

	case protoreflect.DoubleKind:
		if num, ok := value.(float64); ok {
			return protoreflect.ValueOfFloat64(num), nil
		}
		return protoreflect.Value{}, &FieldError{
			Field:    fieldPath,
			Expected: "float64",
			Actual:   fmt.Sprintf("%T", value),
			Message:  "expected 64-bit float",
		}

	case protoreflect.StringKind:
		if str, ok := value.(string); ok {
			return protoreflect.ValueOfString(str), nil
		}
		return protoreflect.Value{}, &FieldError{
			Field:    fieldPath,
			Expected: "string",
			Actual:   fmt.Sprintf("%T", value),
			Message:  "expected string value",
		}

	case protoreflect.BytesKind:
		if str, ok := value.(string); ok {
			return protoreflect.ValueOfBytes([]byte(str)), nil
		}
		return protoreflect.Value{}, &FieldError{
			Field:    fieldPath,
			Expected: "string (for bytes)",
			Actual:   fmt.Sprintf("%T", value),
			Message:  "expected string value for bytes field",
		}

	case protoreflect.EnumKind:
		return c.convertEnumValue(value, field, fieldPath)

	case protoreflect.MessageKind:
		if jsonObj, ok := value.(map[string]interface{}); ok {
			subMsg := dynamicpb.NewMessage(field.Message())
			var errors []FieldError
			c.populateMessage(subMsg, jsonObj, fieldPath, &errors)
			if len(errors) > 0 {
				return protoreflect.Value{}, &errors[0] // Return first error
			}
			return protoreflect.ValueOfMessage(subMsg), nil
		}
		return protoreflect.Value{}, &FieldError{
			Field:    fieldPath,
			Expected: "object",
			Actual:   fmt.Sprintf("%T", value),
			Message:  "expected object for message field",
		}

	default:
		return protoreflect.Value{}, &FieldError{
			Field:    fieldPath,
			Expected: "supported type",
			Actual:   field.Kind().String(),
			Message:  "unsupported field type",
		}
	}
}

// convertEnumValue converts a JSON value to an enum
func (c *Client) convertEnumValue(value interface{}, field protoreflect.FieldDescriptor,
	fieldPath string) (protoreflect.Value, *FieldError) {

	enumDesc := field.Enum()

	// Try string name first
	if str, ok := value.(string); ok {
		enumVal := enumDesc.Values().ByName(protoreflect.Name(str))
		if enumVal != nil {
			return protoreflect.ValueOfEnum(enumVal.Number()), nil
		}
	}

	// Try numeric value
	if num, ok := value.(float64); ok {
		enumNum := protoreflect.EnumNumber(int32(num))
		if enumDesc.Values().ByNumber(enumNum) != nil {
			return protoreflect.ValueOfEnum(enumNum), nil
		}
	}

	return protoreflect.Value{}, &FieldError{
		Field:    fieldPath,
		Expected: "valid enum name or number",
		Actual:   fmt.Sprintf("%v", value),
		Message:  "invalid enum value",
	}
}

func (c *Client) ProtoToJSON(msg proto.Message) ([]byte, error) {
	marshaler := protojson.MarshalOptions{
		UseProtoNames:   true,  // Use proto field names instead of lowerCamelCase
		EmitUnpopulated: false, // Don't include zero values
		UseEnumNumbers:  false, // Use enum names instead of numbers
		Indent:          "  ",  // Pretty print with indentation
	}
	return marshaler.Marshal(msg)
}

// buildFieldPath constructs a dot-separated field path
func (c *Client) buildFieldPath(basePath, fieldName string) string {
	if basePath == "" {
		return fieldName
	}
	return fmt.Sprintf("%s.%s", basePath, fieldName)
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

// Example usage:
/*
func main() {
	// Load service definitions
	registry := NewServiceRegistry()

	// Load from FileDescriptorSet
	fds := &descriptorpb.FileDescriptorSet{} // loaded from file
	err := registry.LoadFileDescriptorSet(fds)
	if err != nil {
		log.Fatal(err)
	}

	// Create gRPC connection
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Create client
	client := NewClient(conn, registry)

	// Make a call
	requestJSON := []byte(`{"name": "John", "age": 30}`)
	responseJSON, err := client.Call(
		context.Background(),
		"com.example.UserService",
		"GetUser",
		requestJSON,
		&CallOptions{},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response:", string(responseJSON))
}
*/
