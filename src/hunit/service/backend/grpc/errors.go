package grpc

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var grpcErrorCodes = map[string]codes.Code{
	"OK":                  codes.OK,
	"CANCELLED":           codes.Canceled, //sic
	"UNKNOWN":             codes.Unknown,
	"INVALID_ARGUMENT":    codes.InvalidArgument,
	"DEADLINE_EXCEEDED":   codes.DeadlineExceeded,
	"NOT_FOUND":           codes.NotFound,
	"ALREADY_EXISTS":      codes.AlreadyExists,
	"PERMISSION_DENIED":   codes.PermissionDenied,
	"RESOURCE_EXHAUSTED":  codes.ResourceExhausted,
	"FAILED_PRECONDITION": codes.FailedPrecondition,
	"ABORTED":             codes.Aborted,
	"OUT_OF_RANGE":        codes.OutOfRange,
	"UNIMPLEMENTED":       codes.Unimplemented,
	"INTERNAL":            codes.Internal,
	"UNAVAILABLE":         codes.Unavailable,
	"DATA_LOSS":           codes.DataLoss,
	"UNAUTHENTICATED":     codes.Unauthenticated,
}

func parseErrorCode(c any, d codes.Code) codes.Code {
	switch x := c.(type) {
	case int:
		return codes.Code(x)
	case float64:
		return codes.Code(x)
	case string:
		v, ok := grpcErrorCodes[strings.ToUpper(x)]
		if ok {
			return v
		}
	}
	return d
}

func grpcErrf(c codes.Code, f string, a ...any) error {
	return status.Error(c, "instaunit: "+fmt.Sprintf(f, a...))
}
