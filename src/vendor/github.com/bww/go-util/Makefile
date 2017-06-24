
TEST_PKGS = ./env \
						./text \
						./rand \
						./qname \
						./uuid \
						./slug

.PHONY: all test

all: test

test:
	go test -test.v $(TEST_PKGS)
