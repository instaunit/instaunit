package emit

import (
  "fmt"
)

/**
 * Documentation format type
 */
type Doctype uint32

const (
  DoctypeMarkdown Doctype = iota
  DoctypeConfluence
  DoctypeInvalid
)

var doctypeNames = []string{
  "markdown",
  "confluence",
  "<invalid>",
}

/**
 * Stringer
 */
func (c Doctype) String() string {
  if c < 0 || c >= DoctypeInvalid {
    return "<invalid>"
  }else{
    return doctypeNames[int(c)]
  }
}

/**
 * Marshal
 */
func (d Doctype) MarshalYAML() (interface{}, error) {
  return d.String(), nil
}

/**
 * Unmarshal
 */
func (d *Doctype) UnmarshalYAML(unmarshal func(interface{}) error) error {
  var s string
  err := unmarshal(&s)
  if err != nil {
    return err
  }
  switch s {
    case "markdown":
      *d = DoctypeMarkdown
    case "confluence":
      *d = DoctypeConfluence
    default:
      return fmt.Errorf("Unsupported documentation type: %v", s)
  }
  return nil
}
