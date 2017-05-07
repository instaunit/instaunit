
export GOPATH := $(GOPATH):$(PWD)

.PHONY: all deps test

all: test

deps:

test:
	go test -test.v ./rand ./qname ./uuid
