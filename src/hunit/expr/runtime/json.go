package runtime

import (
	"encoding/json"
	"fmt"
)

// json libs
type stdJSON struct{}

func (s stdJSON) Marshal(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s stdJSON) Unmarshal(v interface{}) (interface{}, error) {
	var data []byte
	switch c := v.(type) {
	case []byte:
		data = c
	case string:
		data = []byte(c)
	default:
		return nil, fmt.Errorf("Unsupported type: %T", v)
	}
	var obj interface{}
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
