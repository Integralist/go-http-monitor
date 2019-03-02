commit := $(shell git rev-parse --short HEAD)
ldflags := -ldflags "-X main.version=$(commit)"
application := cmd/httpmon/main.go
binary := dist/httpmon

# user can pass cli flags (as a string) to Makefile's `run` target.
# this flags will be passed directly onto the golang program.
#
# example: make run flags="-h"

test:
	go test -v -failfast ./...

run:
	@go run $(ldflags) $(application) ${flags}

build:
	go build $(ldflags) -o $(binary) $(application)

clean:
	@rm $(binary) &> /dev/null || true
