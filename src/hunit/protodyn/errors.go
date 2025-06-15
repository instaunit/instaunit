package protodyn

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCError struct {
	Status  codes.Code
	Message string
	Cause   error
}

func newGRPCError(status codes.Code, err error, msg string) GRPCError {
	return GRPCError{
		Status:  status,
		Message: msg,
		Cause:   err,
	}
}

func (e GRPCError) Error() string {
	return e.Message
}

var (
	ErrNotFound         = errors.New("service or method not found")
	ErrBadRequest       = errors.New("invalid request")
	ErrUnauthenticated  = errors.New("authentication required")
	ErrPermissionDenied = errors.New("permission denied")
	ErrUnavailable      = errors.New("service unavailable")
	ErrTimeout          = errors.New("request timeout")
	ErrGRPCUnknown      = errors.New("gRPC error")
)

// grpcErr converts gRPC errors to useful ones
func grpcErr(err error) error {
	if st, ok := status.FromError(err); ok {
		switch c := st.Code(); c {
		case codes.NotFound:
			return newGRPCError(c, ErrNotFound, st.Message())
		case codes.InvalidArgument:
			return newGRPCError(c, ErrBadRequest, st.Message())
		case codes.Unauthenticated:
			return newGRPCError(c, ErrUnauthenticated, st.Message())
		case codes.PermissionDenied:
			return newGRPCError(c, ErrPermissionDenied, st.Message())
		case codes.Unavailable:
			return newGRPCError(c, ErrUnavailable, st.Message())
		case codes.DeadlineExceeded:
			return newGRPCError(c, ErrTimeout, st.Message())
		default:
			return newGRPCError(c, fmt.Errorf("gRPC error: %v", c), st.Message())
		}
	}
	return err
}
