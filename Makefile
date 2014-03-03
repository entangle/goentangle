SOURCE := $(wildcard *.go)
LIBRARIES := \
	github.com/vmihailenco/msgpack \
	code.google.com/p/snappy-go/snappy
LIBRARIES_DIRS := $(addprefix src/, $(LIBRARIES))

export GOPATH=$(shell pwd)

all:

$(LIBRARIES_DIRS):
	@go get $(@:src/%=%)

test: $(LIBRARIES_DIRS)
	@go test -v .

format:
	@gofmt -l -w $(SOURCE)

clean:
	@rm -rf bin pkg dist

.PHONY: test clean
