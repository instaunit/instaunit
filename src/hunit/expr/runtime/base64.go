package runtime

import (
	"encoding/base64"
	"fmt"
)

// base64 libs
type stdBase64 struct{}

func (s stdBase64) Encode(v interface{}) (string, error) {
	switch c := v.(type) {
	case []byte:
		return base64.StdEncoding.EncodeToString(c), nil
	case string:
		return base64.StdEncoding.EncodeToString([]byte(c)), nil
	default:
		return "", fmt.Errorf("Unsupported type: %T", v)
	}
}

func (s stdBase64) Decode(v interface{}) (string, error) {
	var d []byte
	var err error
	switch c := v.(type) {
	case []byte:
		d, err = base64.StdEncoding.DecodeString(string(c))
	case string:
		d, err = base64.StdEncoding.DecodeString(c)
	default:
		err = fmt.Errorf("Unsupported type: %T", v)
	}
	if err != nil {
		return "", err
	}
	return string(d), nil
}
