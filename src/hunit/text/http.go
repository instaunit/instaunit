package text

import (
	"io"
	"net/http"
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
