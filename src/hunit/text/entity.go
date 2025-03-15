package text

import (
	"encoding/json"
	"fmt"

	"github.com/instaunit/instaunit/hunit/httputil"
	"github.com/instaunit/instaunit/hunit/httputil/mimetype"
)

var (
	ErrUnsupportedContentType = fmt.Errorf("Unsupported content type")
)

// formats entities
type entityFormatter func([]byte) ([]byte, error)

// content types to formatters
var stdEntityFormatters = map[string]entityFormatter{
	mimetype.JSON: jsonEntityFormatter,
}

/**
 * Format an entity, if possible
 */
func FormatEntity(entity []byte, contentType string) ([]byte, error) {
	for k, f := range stdEntityFormatters {
		if httputil.MatchesContentType(k, contentType) {
			return f(entity)
		}
	}
	return nil, ErrUnsupportedContentType
}

/**
 * Format json
 */
func jsonEntityFormatter(entity []byte) ([]byte, error) {
	var v interface{}
	err := json.Unmarshal(entity, &v)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(v, "", "  ")
}
