package test

import (
	"fmt"
)

// IOMode types
type IOMode uint32

const (
	IOModeSync IOMode = iota
	IOModeAsync
)

var iomodeNames = []string{
	"blocking",
	"async",
}

// Stringer
func (c IOMode) String() string {
	if c < 0 || c > IOModeAsync {
		return "<invalid>"
	} else {
		return iomodeNames[int(c)]
	}
}

// Marshal
func (c IOMode) MarshalYAML() (interface{}, error) {
	return c.String(), nil
}

// Unmarshal
func (c *IOMode) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}
	switch s {
	case "blocking", "sync", "":
		*c = IOModeSync
	case "async":
		*c = IOModeAsync
	default:
		return fmt.Errorf("Unsupported I/O mode: %v", s)
	}
	return nil
}
