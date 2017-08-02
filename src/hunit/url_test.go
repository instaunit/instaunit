package hunit

import (
  "fmt"
  "testing"
  "github.com/stretchr/testify/assert"
)

/**
 * Test UUID resemblance
 */
func TestURLMergeQuery(t *testing.T) {
  r, err := mergeQueryParams("file", map[string]string{"a":"b", "c":"d"})
  if assert.Nil(t, err, fmt.Sprintf("%v", err)) {
    assert.Equal(t, "file?a=b&c=d", r)
  }
  r, err = mergeQueryParams("http://host/", map[string]string{"a":"b", "c":"d"})
  if assert.Nil(t, err, fmt.Sprintf("%v", err)) {
    assert.Equal(t, "http://host/?a=b&c=d", r)
  }
  r, err = mergeQueryParams("http://host/file?x=y", map[string]string{"a":"b", "c":"d"})
  if assert.Nil(t, err, fmt.Sprintf("%v", err)) {
    assert.Equal(t, "http://host/file?a=b&c=d&x=y", r)
  }
  r, err = mergeQueryParams("http://host/file?a=x", map[string]string{"a":"b", "c":"d"})
  if assert.Nil(t, err, fmt.Sprintf("%v", err)) {
    assert.Equal(t, "http://host/file?a=x&a=b&c=d", r)
  }
  r, err = mergeQueryParams("http://host/file?a=b", map[string]string{"a":"b", "c":"d"})
  if assert.Nil(t, err, fmt.Sprintf("%v", err)) {
    assert.Equal(t, "http://host/file?a=b&a=b&c=d", r)
  }
}

