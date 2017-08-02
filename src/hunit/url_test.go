package hunit

import (
  "fmt"
  "testing"
  "github.com/stretchr/testify/assert"
)

/**
 * Test URL resemblance
 */
func TestAbsoluteURL(t *testing.T) {
  assert.Equal(t, true, isAbsoluteURL("http://host"))
  assert.Equal(t, true, isAbsoluteURL("http://"))
  assert.Equal(t, true, isAbsoluteURL("a://"))
  assert.Equal(t, false, isAbsoluteURL("file"))
  assert.Equal(t, false, isAbsoluteURL("://"))
  assert.Equal(t, false, isAbsoluteURL("a//"))
  assert.Equal(t, false, isAbsoluteURL("a:/"))
}

/**
 * Test merge query strings
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

