package slug

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

/**
 * Test interpolate
 */
func TestGithubSlug(t *testing.T) {
  var s string
  var c map[string]int
  
  // avoid collisions
  s, c = Github("The Slug!", c)
  assert.Equal(t, "the-slug", s)
  s, c = Github("The Slug!", c)
  assert.Equal(t, "the-slug-1", s)
  s, c = Github("The Slug!", c)
  assert.Equal(t, "the-slug-2", s)
  
  s, c = Github("The Slug! 1", c)
  assert.Equal(t, "the-slug-1-1", s)
  
}
