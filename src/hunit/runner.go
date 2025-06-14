package hunit

import (
	"bytes"
	stdcontext "context"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bww/go-util/v1/debug"
	textutil "github.com/bww/go-util/v1/text"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/instaunit/instaunit/hunit/assert"
	"github.com/instaunit/instaunit/hunit/entity"
	"github.com/instaunit/instaunit/hunit/expr"
	"github.com/instaunit/instaunit/hunit/httputil/mimetype"
	"github.com/instaunit/instaunit/hunit/protodyn"
	"github.com/instaunit/instaunit/hunit/runtime"
	"github.com/instaunit/instaunit/hunit/testcase"
	"github.com/instaunit/instaunit/hunit/text"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// Run a test case
func RunTest(suite *testcase.Suite, tcase testcase.Case, context runtime.Context) (*Result, FutureResult, expr.Variables, error) {
	start := time.Now()

	// wait if we need to
	if tcase.Wait > 0 {
		<-time.After(tcase.Wait)
	}

	// start with an unevaluated result
	result := &Result{Name: formatUnknownName(tcase), Success: true, Case: tcase}
	defer func() {
		result.Runtime = time.Since(start)
	}()

	// process variables first, they can be referenced by this case, itself
	locals := make(expr.Variables)
	for k, e := range tcase.Vars {
		v := textutil.Stringer(e)
		r, err := context.Interpolate(v)
		if err != nil {
			return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, nil, nil
		}
		locals[k] = r
	}

	// test case variables
	vars := expr.Variables{
		"test": tcase,
		"vars": locals,
	}
	context.AddVars(vars)

	if tcase.RPC != nil {
		return runGRPC(suite, tcase, vars, result, context)
	} else {
		return runREST(suite, tcase, vars, result, context)
	}
}

// Run a REST test case
func runREST(suite *testcase.Suite, tcase testcase.Case, vars expr.Variables, result *Result, context runtime.Context) (*Result, FutureResult, expr.Variables, error) {
	var vdef expr.Variables

	// incrementally update the name as we evaluate it
	result.Name = formatRESTName(tcase, tcase.Request.Method, tcase.Request.URL)

	// update the method
	method, err := context.Interpolate(tcase.Request.Method)
	if err != nil {
		return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
	} else if method == "" {
		return nil, nil, nil, fmt.Errorf("Request requires a method (set 'method')")
	}

	// update the url
	url, err := context.Interpolate(tcase.Request.URL)
	if err != nil {
		return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
	}
	if !isAbsoluteURL(url) {
		url = joinPath(context.BaseURL, url)
	}

	// incrementally update the name as we evaluate it
	result.Name = formatRESTName(tcase, method, url)

	if len(tcase.Request.Params) > 0 {
		url, err = mergeQueryParams(url, tcase.Request.Params, context)
		if err != nil {
			return nil, nil, vars, fmt.Errorf("Test case declared on line %d: %v", tcase.Source.Line, err)
		}
	}

	// incrementally update the name as we evaluate it
	result.Name = formatRESTName(tcase, method, url)

	url, err = context.Interpolate(url)
	if err != nil {
		return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
	} else if url == "" {
		return nil, nil, nil, fmt.Errorf("Request requires a URL (set 'url')")
	}

	// incrementally update the name as we evaluate it
	result.Name = formatRESTName(tcase, method, url)

	header := make(http.Header)
	if context.Headers != nil {
		for k, v := range context.Headers {
			k, err = context.Interpolate(k)
			if err != nil {
				return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
			}
			v, err = context.Interpolate(v)
			if err != nil {
				return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
			}
			header.Add(k, v)
		}
	}

	if tcase.Request.Headers != nil {
		for k, v := range tcase.Request.Headers {
			k, err = context.Interpolate(k)
			if err != nil {
				return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
			}
			v, err = context.Interpolate(v)
			if err != nil {
				return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
			}
			header.Add(k, v)
		}
	}

	if tcase.Request.BasicAuth != nil {
		u, err := context.Interpolate(tcase.Request.BasicAuth.Username)
		if err != nil {
			return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
		}
		p, err := context.Interpolate(tcase.Request.BasicAuth.Password)
		if err != nil {
			return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
		}
		header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(u+":"+p))))
	}

	// if we expect a stream, we must set it up and return our future result here
	if tcase.Stream != nil {
		messages := tcase.Stream.Messages
		if len(messages) < 1 {
			return result.Error(fmt.Errorf("No messages are exchanged over websocket.")), nil, vars, nil
		}

		dialer := websocket.Dialer{
			NetDial: func(n, a string) (net.Conn, error) {
				return net.DialTimeout(n, a, time.Second*3)
			},
		}
		url, err := urlWithScheme("ws", url)
		if err != nil {
			return result.Error(fmt.Errorf("Could not upgrade URL scheme: %w", err)), nil, vars, nil
		}
		conn, _, err := dialer.Dial(url, header)
		if err != nil {
			return result.Error(fmt.Errorf("Could not dial websocket: %w", err)), nil, vars, nil
		}

		monitor := NewStreamMonitor(url, context, conn, messages)
		err = monitor.Run(result)
		if err != nil {
			return nil, nil, nil, err
		}

		// websocket variables
		vdef = expr.Variables{
			"url": url,
		}
		vars["websocket"] = vdef
		context.AddVars(expr.Variables{
			"websocket": vdef,
		})

		// depending on the I/O mode, resolve or return a future
		switch m := tcase.Stream.Mode; m {
		case testcase.IOModeAsync:
			return nil, monitor, vars, nil
		case testcase.IOModeSync:
			r, err := monitor.Finish(time.Time{})
			if err != nil {
				return result.Error(fmt.Errorf("Could not finish I/O: %w", err)), nil, vars, nil
			}
			return r, nil, vars, nil
		default:
			return nil, nil, nil, fmt.Errorf("No such I/O mode: %v", m)
		}
	}

	var reqdata string
	var ereader io.Reader
	if tcase.Request.Entity != "" {
		reqdata, err = context.Interpolate(tcase.Request.Entity)
		if err != nil {
			return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
		} else {
			ereader = bytes.NewBuffer([]byte(reqdata))
		}
	}
	if reqdata != "" {
		header.Add("Content-Length", strconv.FormatInt(int64(len(reqdata)), 10))
	}

	req, err := http.NewRequest(method, url, ereader)
	if err != nil {
		return nil, nil, nil, err
	}

	req.Header = header

	if tcase.Request.Cookies != nil {
		for k, v := range tcase.Request.Cookies {
			k, err = context.Interpolate(k)
			if err != nil {
				return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
			}
			v, err = context.Interpolate(v)
			if err != nil {
				return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
			}
			req.AddCookie(&http.Cookie{Name: k, Value: v})
		}
	}

	reqbuf := &bytes.Buffer{}
	err = text.WriteRequest(reqbuf, req, reqdata)
	if err != nil {
		return result.Error(fmt.Errorf("Could not write: %w", err)), nil, vars, nil
	} else {
		result.Reqdata = reqbuf.Bytes()
	}

	rsp, err := context.Client.Do(req)
	if rsp != nil && rsp.Body != nil {
		defer rsp.Body.Close()
	}
	if err != nil {
		return result.Error(fmt.Errorf("Could not read response body: %w", err)), nil, vars, nil
	}

	// check the response status
	if tcase.Response.Status == 0 { // if the status is not explicitly defined we assume 200/OK is expected
		result.AssertEqual(http.StatusOK, rsp.StatusCode, "Unexpected status code (default)")
	} else {
		result.AssertEqual(tcase.Response.Status, rsp.StatusCode, "Unexpected status code")
	}

	// note the content type; we prefer the explicit format over the content type
	contentType, err := context.Interpolate(textutil.Coalesce(
		tcase.Response.Format,
		strings.ToLower(rsp.Header.Get("Content-Type")),
	))
	if err != nil {
		return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
	}

	// check response headers, if necessary
	if headers := tcase.Response.Headers; headers != nil {
		for k, v := range headers {
			k, err = context.Interpolate(k)
			if err != nil {
				return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
			}
			v, err = context.Interpolate(v)
			if err != nil {
				return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
			}
			result.AssertEqual(v, rsp.Header.Get(k), "Headers do not match: %v", k)
		}
	}

	// Transform the response, if necessary, first by applying suite-level
	// transforms, then request-level transforms. We do not currently make an
	// attempt to avoid applying the same transform repeatedly, either at
	// different levels or the same level.
	var xfspecs []testcase.Transform
	if specs := suite.Transform.Response; specs != nil {
		xfspecs = append(xfspecs, specs...)
	}
	if specs := tcase.Response.Transforms; specs != nil {
		xfspecs = append(xfspecs, specs...)
	}
	for _, spec := range xfspecs {
		rsp, err = spec.TransformResponse(rsp)
		if err != nil {
			return result.Error(fmt.Errorf("Could not transform response: %w", err)), nil, vars, nil
		}
	}

	// handle the response entity
	var rspdata []byte
	if rsp.Body != nil {
		rspdata, err = io.ReadAll(rsp.Body)
		if err != nil {
			result.Error(fmt.Errorf("Could not read response body: %w", err))
		}
	}

	// parse response entity if necessry
	var rspvalue interface{} = rspdata
	if tcase.Response.Comparison == testcase.CompareSemantic {
		rspvalue, err = entity.Unmarshal(contentType, rspdata)
		if err != nil {
			return result.Error(fmt.Errorf("Could not unmarshal entity: %w", err)), nil, vars, nil
		}
	} else if tcase.Id != "" { // attempt it but don't produce an error if we fail
		val, err := entity.Unmarshal(contentType, rspdata)
		if err == nil {
			rspvalue = val
		}
	}

	// check response entity, if necessary
	if entity := tcase.Response.Entity; entity != "" {
		entity, err = context.Interpolate(entity)
		if err != nil {
			return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
		}
		if len(rspdata) == 0 {
			result.AssertEqual(entity, "", "Entities do not match")
		} else if err = entitiesEqual(context, tcase.Response.Comparison, contentType, []byte(entity), rspvalue); err != nil {
			result.Error(fmt.Errorf("Could not compare entities: %w", err))
		}
	}

	rspbuf := &bytes.Buffer{}
	err = text.WriteResponse(rspbuf, rsp, rspdata)
	if err != nil {
		return result.Error(fmt.Errorf("Could not write: %w", err)), nil, vars, nil
	} else {
		result.Rspdata = rspbuf.Bytes()
	}

	// response variables
	vdef = expr.Variables{
		"headers": flattenHeader(rsp.Header),
		"cookies": flattenCookies(rsp.Cookies()),
		"entity":  rspdata,
		"value":   rspvalue,
		"status":  rsp.StatusCode,
	}
	vars["response"] = vdef
	context.AddVars(expr.Variables{
		"response": vdef,
	})

	// update request with final context
	result.Context = context

	// update test case dynamic post-fields with response
	if r := text.Coalesce(tcase.Route.Id, suite.Route.Id); r != "" {
		r, err = context.Interpolate(r)
		if err != nil {
			return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
		}
		tcase.Route.Id = r
	}
	if r := text.Coalesce(tcase.Route.Path, suite.Route.Path); r != "" {
		r, err = context.Interpolate(r)
		if err != nil {
			return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
		}
		tcase.Route.Path = r
	}

	// assertions
	if assert := tcase.Response.Assert; assert != nil {
		ok, err := assert.Bool(context.Variables)
		if err != nil {
			b := &strings.Builder{}
			debug.Dumpf(b, context.Variables)
			return nil, nil, nil, fmt.Errorf("Could not evaluate assertion: %v\n%s", err, b.String())
		}
		if !ok {
			result.Error(&ScriptError{"Script assertion failed", true, ok, assert})
		}
	}

	// generate documentation if necessary
	if tcase.Documented() && len(context.Gendoc) > 0 {
		for _, e := range context.Gendoc {
			l, err := tcase.Interpolate(context.Variables)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("Could not generate documentation: %w", err)
			}
			err = e.Case(suite, l, req, reqdata, rsp, rspdata)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("Could not generate documentation: %w", err)
			}
		}
	}

	return result, nil, vars, nil
}

// Run a gRPC test case
func runGRPC(suite *testcase.Suite, tcase testcase.Case, vars expr.Variables, result *Result, context runtime.Context) (*Result, FutureResult, expr.Variables, error) {
	if context.Service == nil {
		return nil, nil, nil, fmt.Errorf("gRPC test case requires at least one gRPC service is defined")
	}

	// incrementally update the name as we evaluate it
	result.Name = formatGRPCName(tcase)

	// attempt to connect to the service (we connect for each request, which isn't
	// performant, but is more suitable for our needs here).
	conn, err := grpc.Dial(tcase.Request.URL, grpc.WithInsecure())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Could not connect to service: %s", tcase.Request.URL)
	}
	defer conn.Close()

	// create a client for the specified gRPC service
	client := protodyn.NewClient(conn, context.Service)
	cxt := stdcontext.Background()

	// create an invocation for the RPC call
	inv, err := client.Invocation(cxt, tcase.RPC.Service, tcase.RPC.Method, &protodyn.CallOptions{})
	if err != nil {
		return result.Error(fmt.Errorf("Could not find gRPC method: %w", err)), nil, vars, nil
	}

	// setup the request entity
	var (
		reqmsg  *dynamicpb.Message
		reqdata []byte
	)
	if ent := tcase.Request.Entity; ent != "" {
		reqtype, err := context.Interpolate(tcase.Request.Format)
		if err != nil {
			return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
		}
		reqmsg, reqdata, err = unmarshalGRPCRequest(context, inv, reqtype, ent)
		if err != nil {
			return result.Error(err), nil, vars, nil
		}
	}

	// update the request data in the result
	result.Reqdata = reqdata

	// perform the gRPC request
	rspmsg, err := client.Invoke(cxt, inv, reqmsg)
	if err != nil {
		return result.Error(fmt.Errorf("Could not call gRPC method: %w", err)), nil, vars, nil
	}

	// decode the response entity to JSON
	rspdata, err := protodyn.MarshalJSON(rspmsg)
	if err != nil {
		return result.Error(fmt.Errorf("Could not convert gRPC response: %w", err)), nil, vars, nil
	}
	// update the request data in the result
	result.Rspdata = rspdata
	// unmarshal it to the intermediate format
	rspvalue, err := entity.Unmarshal(mimetype.JSON, rspdata)
	if err != nil {
		return result.Error(fmt.Errorf("Could not decode gRPC response: %w", err)), nil, vars, nil
	}

	// check response entity, if necessary
	if expectent := tcase.Response.Entity; expectent != "" {
		expecttype, err := context.Interpolate(tcase.Response.Format)
		if err != nil {
			return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
		}
		expect, expectdata, err := unmarshalGRPCResponse(context, inv, expecttype, expectent)
		if err != nil {
			return result.Error(err), nil, vars, nil
		}
		// If the expected response type is defined in JSON, this is the lowest
		// common denominator, so we use it to compare messages. Only when both
		// messages are defined as protos do we compare them directly.
		switch expecttype {
		case mimetype.Protobuf:
			if !proto.Equal(expect, rspmsg) {
				result.Error(&assert.AssertionError{
					Expect:  expect,
					Actual:  rspmsg,
					Message: "Entities are not equal (protobuf)",
				})
			}
		case mimetype.JSON:
			fallthrough
		default:
			expect, err := entity.Unmarshal(mimetype.JSON, expectdata)
			if err != nil {
				return result.Error(fmt.Errorf("Could not unmarshal expected response: %w", err)), nil, vars, nil
			}
			if len(rspdata) == 0 {
				result.AssertEqual(expect, "", "Entities do not match")
			} else if !entity.SemanticEqual(expect, rspvalue) {
				result.Error(&assert.AssertionError{
					Expect:  expect,
					Actual:  rspvalue,
					Message: "Entities are not equal",
				})
			}
		}
	}

	vdef := expr.Variables{
		"entity": rspdata,
		"value":  rspvalue,
	}
	vars["response"] = vdef
	context.AddVars(expr.Variables{})

	// update result with final context
	result.Context = context

	return result, nil, vars, nil
}

func unmarshalGRPCRequest(context runtime.Context, inv protodyn.Invocation, ctype, edata string) (*dynamicpb.Message, []byte, error) {
	return unmarshalGRPC(context, inv, inv.Method.Input(), ctype, edata)
}

func unmarshalGRPCResponse(context runtime.Context, inv protodyn.Invocation, ctype, edata string) (*dynamicpb.Message, []byte, error) {
	return unmarshalGRPC(context, inv, inv.Method.Output(), ctype, edata)
}

func unmarshalGRPC(context runtime.Context, inv protodyn.Invocation, md protoreflect.MessageDescriptor, ctype, edata string) (*dynamicpb.Message, []byte, error) {
	var (
		entdata []byte
		entmsg  *dynamicpb.Message
		err     error
	)
	switch ctype {
	// if the entity is explicitly a protobuf message, we decode it as base64
	// and use it directly...
	case mimetype.Protobuf:
		entdata, err = base64.StdEncoding.DecodeString(edata)
		if err != nil {
			return nil, nil, fmt.Errorf("Could not decode payload: %w", err)
		}
		err = proto.Unmarshal(entdata, entmsg)
		if err != nil {
			return nil, nil, fmt.Errorf("Could not unmarshal protobuf message: %w", err)
		}

	// ...otherwise, it is assumed to be JSON data which we will marshal to the
	// corresponding protobuf message.
	case mimetype.JSON:
		fallthrough
	default:
		reqstr, err := context.Interpolate(string(edata))
		if err != nil {
			return nil, nil, fmt.Errorf("Could not interpolate: %w", err)
		}
		entdata = []byte(reqstr)
		entmsg, err = inv.MessageFromJSON(md, entdata)
		if err != nil {
			return nil, nil, fmt.Errorf("Could not unmarshal message from JSON: %w", err)
		}
	}
	return entmsg, entdata, nil
}
