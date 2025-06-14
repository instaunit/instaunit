package protodyn

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// default JSON marshaler
var jsonMarshaler = protojson.MarshalOptions{
	UseProtoNames:   true,  // Use proto field names instead of lowerCamelCase
	EmitUnpopulated: false, // Don't include zero values
	UseEnumNumbers:  false, // Use enum names instead of numbers
	Indent:          "  ",  // Pretty print with indentation
}

func MarshalJSON(msg proto.Message) ([]byte, error) {
	return jsonMarshaler.Marshal(msg)
}

// default JSON unmarshaler
var jsonUnmarshaler = protojson.UnmarshalOptions{}

func UnmarshalJSON(data []byte, msg proto.Message) error {
	return jsonUnmarshaler.Unmarshal(data, msg)
}
