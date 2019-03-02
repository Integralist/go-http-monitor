commit := $(shell git rev-parse --short HEAD)
ldflags := -ldflags "-X main.version=$(commit)"
application := cmd/httpmon/main.go
binary := dist/httpmon

test:
	go test -v -failfast ./...

run:
	@go run $(ldflags) $(application)

build:
	go build $(ldflags) -o $(binary) $(application)

clean:
	@rm $(binary) &> /dev/null || true
