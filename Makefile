APPNAME ?= chat-test

# used by `test` target
export REPORTS_DIR=./reports
# used by lint target
export GOMETALINTER_CONFIG=.gometalinter.json
export GOLANGCILINT_VERSION=v1.23.8

build: clean
	mkdir -p build
	GOOS=$(GOOS) GOARCH=$(GOARCH) APPNAME=$(APPNAME) ./scripts/build

run: build
	./build/${APPNAME}

test:
	./scripts/unit-test

test-report:
	./scripts/show-tests

lint:
	./scripts/lint

clean:
	APPNAME=$(APPNAME) ./scripts/clean

.PHONY: build run test test-report lint clean
