commit := $(shell git rev-parse --short HEAD)
ldflags := -ldflags "-X main.version=$(commit)"
application := cmd/httpmon/main.go
binary := dist/httpmon

# user can pass cli flags (as a string) to Makefile's `run` target.
# this flags will be passed directly onto the golang program.
#
# examples:
# 	make run flags="-h"
# 	make run flags="-stats 2"
# 	make run flags="-stats 2 -threshold 5"

test:
	go test -v -failfast ./...

run: clean
	@go run $(ldflags) $(application) ${flags}

build:
	go build $(ldflags) -o $(binary) $(application)

clean:
	@echo '127.0.0.1 - integralist [02/March/2019:09:00:00 +0000] "GET /foo HTTP/1.1" 200 123' > access.log
	@go mod tidy
	@rm $(binary) &> /dev/null || true
