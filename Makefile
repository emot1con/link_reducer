build:
	@go build -o bin/go_link_reducer cmd/main.go

test:
	@go test -v ./..

run: build
	@./bin/go_link_reducer