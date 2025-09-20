package grpc

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"

	yaml "gopkg.in/yaml.v3"
)

type RemoteProcedure struct {
	Service string `yaml:"service"`
	Method  string `yaml:"method"`
}

// A request
type Request struct {
	Headers map[string]string `yaml:"headers"`
	Format  string            `yaml:"format"` // the format of entity (only 'application/json' is currently supported)
	Entity  string            `yaml:"entity"`
}

// A response
type Response struct {
	Status  int               `yaml:"status"`
	Headers map[string]string `yaml:"headers"`
	Cookies map[string]string `yaml:"cookies"`
	Entity  string            `yaml:"entity"`
}

// An endpoint
type Endpoint struct {
	Wait     time.Duration    `yaml:"wait"`
	RPC      *RemoteProcedure `yaml:"grpc"`
	Request  *Request         `yaml:"request"`
	Response *Response        `yaml:"response"`
}

// A test suite
type Suite struct {
	Protos    []string   `yaml:"protos"`
	Endpoints []Endpoint `yaml:"service"`

	epByMethodOnce sync.Once
	epByMethod     map[string][]Endpoint
}

// Load a test suite
func LoadSuite(src io.ReadCloser) (*Suite, error) {
	data, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	suite := &Suite{}
	err = unmarshal(data, suite)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal gRPC service: %w", err)
	}

	return suite, nil
}

func (s *Suite) MatchEndpoint(mname string, reqmsg *dynamicpb.Message, matchers ...RequestMatcher) (Endpoint, bool) {
	s.epByMethodOnce.Do(func() {
		m := make(map[string][]Endpoint)
		for _, e := range s.Endpoints {
			if rpc := e.RPC; rpc != nil {
				name := methodName(rpc.Service, rpc.Method)
				m[name] = append(m[name], e)
			}
		}
		s.epByMethod = m
	})
	epoints, ok := s.epByMethod[mname]
	if !ok || len(epoints) == 0 {
		return Endpoint{}, false
	}
outer:
	for _, e := range epoints {
		// compare the declared request (if we have one defined) against the actual request
		if req := e.Request; req != nil {
			if ent := req.Entity; ent != "" {
				// Create request message to receive the incoming data
				chkmsg := dynamicpb.NewMessage(reqmsg.Descriptor())
				// Use protojson to unmarshal into the dynamic message
				err := protojson.Unmarshal([]byte(ent), chkmsg)
				if err != nil {
					logf("Could not unmarshal JSON response to protobuf to match gRPC request: %v", err)
					continue outer
				}
				// match the entity
				if !leftEqual(chkmsg, reqmsg) {
					continue outer
				}
			}
		}
		// apply additional user-specified matching criteria
		if !RequestMatchers(matchers).MatchesRequest(e, reqmsg) {
			continue outer
		}
		// if we've reached the end, we've fully matched
		return e, true
	}
	return Endpoint{}, ok
}

func methodName(s, m string) string {
	return fmt.Sprintf("%s/%s", s, m)
}

func unmarshal(data []byte, dest interface{}) error {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	return dec.Decode(dest)
}

// leftEqual compares all the fields in left to the corresponding fields in
// right. Only fields present in left are considered; if a field is present in
// right but not left, it is ignored.
func leftEqual(left, right *dynamicpb.Message) bool {
	match := true // matches unless we have a mismatch
	left.Range(func(fd protoreflect.FieldDescriptor, lval protoreflect.Value) bool {
		if !right.Has(fd) {
			match = false
			return false
		}
		rval := right.Get(fd)
		if !lval.Equal(rval) {
			match = false
			return false
		}
		return true
	})
	return match
}
