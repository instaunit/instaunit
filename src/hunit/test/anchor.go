package test

import (
  "fmt"
)

// Anchor style
type AnchorStyle uint32

const (
  AnchorGithub AnchorStyle = iota
  AnchorRails
  AnchorInvalid
)

var anchorStyleNames = []string{
  "github",
  "rails",
  "<invalid>",
}

// Parse an anchor style
func ParseAnchorStyle(s string) AnchorStyle {
  switch s {
    case "github":
      return AnchorGithub
    case "rails":
      return AnchorRails
    default:
      return AnchorInvalid
  }
}

// Stringer
func (c AnchorStyle) String() string {
  if c < 0 || c >= AnchorInvalid {
    return "<invalid>"
  }else{
    return anchorStyleNames[int(c)]
  }
}

// Marshal
func (c AnchorStyle) MarshalYAML() (interface{}, error) {
  return c.String(), nil
}

// Unmarshal
func (c *AnchorStyle) UnmarshalYAML(unmarshal func(interface{}) error) error {
  var s string
  err := unmarshal(&s)
  if err != nil {
    return err
  }
  v := ParseAnchorStyle(s)
  if v == AnchorInvalid {
    return fmt.Errorf("Unsupported anchor style: %v", s)
  }
  *c = v
  return nil
}
