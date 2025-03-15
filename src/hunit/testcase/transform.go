package testcase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/instaunit/instaunit/hunit/httputil"
	"github.com/instaunit/instaunit/hunit/httputil/mimetype"
)

var (
	errTransformerNotSupported = errors.New("Transformer not supported")
	errContentTypeNotSupported = errors.New("Content type not supported")
)

type Transform struct {
	Type string `yaml:"type"`
}

func (t Transform) ResponseTransfomer() (ResponseTransformer, error) {
	switch strings.ToLower(t.Type) {
	case "pretty":
		return PrettyTransformer{}, nil
	default:
		return nil, fmt.Errorf("%w: %q", errTransformerNotSupported, t.Type)
	}
}

type ResponseTransformer interface {
	// TransformResponse transforms a request in some way and returns
	// either a shallow copy or the unmodified argument
	TransformResponse(*http.Response) (*http.Response, error)
}

type PrettyTransformer struct{}

func (t PrettyTransformer) TransformResponse(rsp *http.Response) (*http.Response, error) {
	ct := httputil.ContentType(rsp.Header)
	switch {
	case httputil.MatchesContentType(mimetype.JSON, ct):
		return t.transformJSONResponse(rsp)
	default:
		return rsp, fmt.Errorf("%w: %s", errContentTypeNotSupported, ct)
	}
}

func (t PrettyTransformer) transformJSONResponse(rsp *http.Response) (*http.Response, error) {
	if rsp.Body == nil {
		return rsp, nil
	}

	var data json.RawMessage
	err := json.NewDecoder(rsp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("Could not read response: %w", err)
	}

	ptty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("Could not format response: %w", err)
	}

	dup := *rsp
	dup.Body = io.NopCloser(strings.NewReader(string(ptty)))

	return &dup, nil
}
