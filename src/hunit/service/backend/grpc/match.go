package grpc

import "google.golang.org/protobuf/types/dynamicpb"

type RequestMatcher interface {
	MatchesRequest(epoint Endpoint, reqmsg *dynamicpb.Message) bool
}

type RequestMatcherFunc func(epoint Endpoint, reqmsg *dynamicpb.Message) bool

func (m RequestMatcherFunc) MatchesRequest(epoint Endpoint, reqmsg *dynamicpb.Message) bool {
	return m(epoint, reqmsg)
}

type RequestMatchers []RequestMatcher

func (m RequestMatchers) MatchesRequest(epoint Endpoint, reqmsg *dynamicpb.Message) bool {
	for _, e := range m {
		if !e.MatchesRequest(epoint, reqmsg) {
			return false
		}
	}
	return true
}
