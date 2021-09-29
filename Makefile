build:
	go build -o ./bin/ ./...

temperature:
	go run cmd/temperature/main.go