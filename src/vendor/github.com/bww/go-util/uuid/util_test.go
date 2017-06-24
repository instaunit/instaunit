package uuid

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

/**
 * Test UUID resemblance
 */
func TestResemblance(t *testing.T) {
  assert.Equal(t, true, ResemblesUUID("ACE24573-5BD5-4C5F-B143-5E9E17F18BDB"))
  assert.Equal(t, true, ResemblesUUID("ace24573-5bd5-4c5f-b143-5e9e17f18bdb"))
  assert.Equal(t, true, ResemblesUUID("ace24573-5BD5-4C5F-B143-5e9e17F18BDB"))
  assert.Equal(t, false, ResemblesUUID("ace24573-5bd5-4c5f-b143-5e9e17f18bd"))
  assert.Equal(t, false, ResemblesUUID("ace24573-5bd5-4cZ5f-b143-5e9e17f18bd"))
  assert.Equal(t, false, ResemblesUUID(""))
}
