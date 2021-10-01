build:
	go build -o ./bin/ ./...

probe:
	go run cmd/probe/main.go