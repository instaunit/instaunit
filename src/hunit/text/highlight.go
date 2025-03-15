package text

import (
	"github.com/instaunit/instaunit/hunit/httputil"
	"github.com/instaunit/instaunit/hunit/httputil/mimetype"
)

var highlights = map[string]string{
	mimetype.JSON:       "json",
	mimetype.Javascript: "javascript",
}

// Obtain the appropriate highlight for the entity
func EntityHighlight(c string) string {
	for k, v := range highlights {
		if httputil.MatchesContentType(k, c) {
			return v
		}
	}
	return ""
}
