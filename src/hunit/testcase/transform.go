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

type ConditionFailedPolicy string

const (
	ConditionFailedSkip    = "skip" // skip the operation when a condition fails
	ConditionFailedFail    = "fail" // fail with an error when a condition fails
	ConditionFailedDefault = ConditionFailedSkip
)

func (c ConditionFailedPolicy) Check(err error) error {
	if err == nil {
		return nil
	}
	switch c {
	case "": // default case falls through to the default handler
		fallthrough
	case ConditionFailedSkip:
		if errors.Is(err, errContentTypeNotSupported) {
			return nil // content type not supported is ignored when policy is: skip
		}
	}
	return err
}

type TransformCollection struct {
	Request  []Transform `json:"request"`
	Response []Transform `json:"response"`
}

type ResponseTransformer interface {
	// TransformResponse transforms a request in some way and returns
	// either a shallow copy or the unmodified argument
	TransformResponse(*http.Response) (*http.Response, error)
}

type Transform struct {
	Type        string                `yaml:"type"`
	Unsupported ConditionFailedPolicy `json:"unsupported"`
}

func (t Transform) ResponseTransfomer() (ResponseTransformer, error) {
	switch strings.ToLower(t.Type) {
	case "pretty":
		return PrettyTransformer{}, nil
	default:
		return nil, fmt.Errorf("%w: %q", errTransformerNotSupported, t.Type)
	}
}

// TransformResponse implements ResponseTransformer for Transform and is the
// preferred way to apply a Transform since it correctly handles
// condition-failed policies
func (t Transform) TransformResponse(rsp *http.Response) (*http.Response, error) {
	xform, err := t.ResponseTransfomer()
	if err != nil {
		return nil, err
	}
	xconv, err := xform.TransformResponse(rsp)
	if err = t.Unsupported.Check(err); err != nil {
		return nil, err
	}
	return xconv, nil
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
