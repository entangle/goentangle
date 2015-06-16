SOURCE := $(wildcard *.go)
LIBRARIES := \
	github.com/vmihailenco/msgpack \
	github.com/golang/snappy/snappy
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
