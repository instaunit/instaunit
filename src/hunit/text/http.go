package text

import (
  "io"
  "net/http"
)

/**
 * Write a request to the specified output
 */
func WriteRequest(w io.Writer, req *http.Request, entity string) error {
  
  dump := req.Method +" "
  dump += req.URL.Path
  if q := req.URL.RawQuery; q != "" { dump += "?"+ q }
  dump += " "+ req.Proto +"\n"
  
  dump += "Host: "+ req.URL.Host +"\n"
  for k, v := range req.Header {
    dump += k +": "
    for i, e := range v {
      if i > 0 { dump += "," }
      dump += e
    }
    dump += "\n"
  }
  
  dump += "\n"
  if entity != "" {
    dump += entity +"\n"
  }
  
  _, err := w.Write([]byte(dump))
  if err != nil {
    return err
  }
  
  return nil
}

/**
 * Write a response to the specified output
 */
func WriteResponse(wHeaders, wEntity io.Writer, rsp *http.Response, entity []byte) error {
  dump := rsp.Proto +" "+ rsp.Status +"\n"
  
  for k, v := range rsp.Header {
    dump += k +": "
    for i, e := range v {
      if i > 0 { dump += "," }
      dump += e
    }
    dump += "\n"
  }
  
  dump += "\n"
  _, err := wHeaders.Write([]byte(dump))
  if err != nil {
    return err
  }
  if entity != nil {
		e := string(entity) + "\n"
		_, err := wEntity.Write([]byte(e))
		if err != nil {
			return err
		}
  }
  
  return nil
}
