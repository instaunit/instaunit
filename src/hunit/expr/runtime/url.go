package runtime

import (
	"net/url"
)

// url libs
type stdURL struct{}

func (s stdURL) Base(v string) (string, error) {
	u, err := url.Parse(v)
	if err != nil {
		return "", err
	}
	u.RawQuery = ""
	return u.String(), nil
}

func (s stdURL) Query(v string) (string, error) {
	u, err := url.Parse(v)
	if err != nil {
		return "", err
	}
	if u.RawQuery != "" {
		return "?" + u.RawQuery, nil
	} else {
		return "", nil
	}
}
