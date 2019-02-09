package text

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
)

// Write a request to the specified output
func WriteRequest(w io.Writer, req *http.Request, entity string) error {
	var dump string

	if req != nil {
		dump += req.Method + " "
		dump += req.URL.Path
		if q := req.URL.RawQuery; q != "" {
			dump += "?" + q
		}
		dump += " " + req.Proto + "\n"

		dump += "Host: " + req.URL.Host + "\n"
		for k, v := range req.Header {
			dump += k + ": "
			for i, e := range v {
				if i > 0 {
					dump += ","
				}
				dump += e
			}
			dump += "\n"
		}
	}

	if entity != "" {
		dump += "\n"
		dump += entity
	}

	_, err := w.Write([]byte(dump))
	if err != nil {
		return err
	}

	return nil
}

// Write a response to the specified output
func WriteResponse(w io.Writer, rsp *http.Response, entity []byte) error {
	var dump string

	if rsp != nil {
		dump += rsp.Proto + " " + rsp.Status + "\n"

		for k, v := range rsp.Header {
			dump += k + ": "
			for i, e := range v {
				if i > 0 {
					dump += ","
				}
				dump += e
			}
			dump += "\n"
		}
	}

	if entity != nil {
		dump += "\n"
		dump += string(entity)
	}

	_, err := w.Write([]byte(dump))
	if err != nil {
		return err
	}

	return nil
}

// Determine if the provided request has a particular content type
func HasContentType(req *http.Request, t string) bool {
	return MatchesContentType(t, req.Header.Get("Content-Type"))
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
