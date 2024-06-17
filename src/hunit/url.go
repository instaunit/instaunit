package hunit

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/instaunit/instaunit/hunit/runtime"
)

// Scheme matcher
var urlish = regexp.MustCompile("^[a-z]+://")

// Determine if a URL appears to be absolute. It is considered absolute
// if the string begins with a scheme://.
func isAbsoluteURL(u string) bool {
	return urlish.MatchString(u)
}

// Merge query parameters with a map
func mergeQueryParams(u string, p map[string]string, c runtime.Context) (string, error) {
	if len(p) < 1 {
		return u, nil
	}

	v, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	q := v.Query()
	for k, v := range p {
		v, err = interpolateIfRequired(c, v)
		if err != nil {
			return "", fmt.Errorf("Evaluating query parameter %q: %v", k, err)
		}
		q.Add(k, v)
	}

	v.RawQuery = q.Encode()
	return v.String(), nil
}

// Change a URL scheme
func urlWithScheme(s, u string) (string, error) {
	v, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	v.Scheme = s
	return v.String(), nil
}
