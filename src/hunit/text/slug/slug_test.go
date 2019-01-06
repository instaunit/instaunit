package slug

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test interpolate
func TestGithubSlug(t *testing.T) {
	var s string
	var c map[string]int

	// avoid collisions
	c = nil
	s, c = Github("The Slug!", c)
	assert.Equal(t, "the-slug", s)
	s, c = Github("GET This-Slug! Here;    We've even_got_some_underscores.", c)
	assert.Equal(t, "get-thisslug-here-weve-even_got_some_underscores", s)

	// avoid collisions
	c = nil
	s, c = Github("The Slug!", c)
	assert.Equal(t, "the-slug", s)
	s, c = Github("The Slug!", c)
	assert.Equal(t, "the-slug-1", s)
	s, c = Github("The Slug!", c)
	assert.Equal(t, "the-slug-2", s)

	// this is a known edge case which we don't handle
	// keep c
	s, c = Github("The Slug! 1", c)
	assert.Equal(t, "the-slug-1", s) // probably should be the-slug-1-1

}
