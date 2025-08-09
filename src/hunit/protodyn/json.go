package protodyn

import (
	"errors"
	"fmt"

	"github.com/bww/go-util/v1/text"
	"github.com/instaunit/instaunit/hunit/reflect"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var errInvalidMessage = errors.New("Invalid input message")

// default JSON marshaler
var jsonMarshaler = protojson.MarshalOptions{
	UseProtoNames:   true,  // Use proto field names instead of lowerCamelCase
	EmitUnpopulated: false, // Don't include zero values
	UseEnumNumbers:  false, // Use enum names instead of numbers
	Indent:          "  ",  // Pretty print with indentation
}

func MarshalJSON(msg proto.Message) ([]byte, error) {
	if reflect.IsNil(msg) {
		return []byte("null"), nil
	}
	return jsonMarshaler.Marshal(msg)
}

// default JSON unmarshaler
var jsonUnmarshaler = protojson.UnmarshalOptions{}

func UnmarshalJSON(data []byte, msg proto.Message) error {
	if reflect.IsNil(msg) {
		return fmt.Errorf("%w: message is nil", errInvalidMessage)
	}
	err := jsonUnmarshaler.Unmarshal(data, msg)
	if err != nil {
		return fmt.Errorf("%w: in data:\n%s", err, text.Indent(string(data), "> "))
	}
	return nil
}
