package httputil

import (
	"fmt"
	"net/http"
	"path"
	"strings"
)

// Determine if the provided request has a particular content type
func ContentType(hdr http.Header) string {
	return hdr.Get("Content-Type")
}

// Determine if the provided request has a particular content type
func HasContentType(req *http.Request, t string) bool {
	return MatchesContentType(t, ContentType(req.Header))
}

// Determine if the provided request has a particular content type
func MatchesContentType(pattern, contentType string) bool {

	// trim off the parameters following ';' if we have any
	if i := strings.Index(contentType, ";"); i > 0 {
		contentType = contentType[:i]
	}

	// path.Match does glob matching, which is useful it we
	// want to, e.g., test for all image types with `image/*`.
	m, err := path.Match(pattern, contentType)
	if err != nil {
		panic(fmt.Errorf("* * * could not match invalid content-type pattern: %s", pattern))
	}

	return m
}
