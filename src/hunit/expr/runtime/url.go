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

func (s stdURL) AddParam(b, k, v string) (string, error) {
	u, err := url.Parse(b)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Add(k, v)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (s stdURL) SetParam(b, k, v string) (string, error) {
	u, err := url.Parse(b)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Add(k, v)
	u.RawQuery = q.Encode()
	return u.String(), nil
}
