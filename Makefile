build:
	go build -o ./bin/probe/ ./cmd/probe/probe.go
	go build -o ./bin/subscriber/ ./cmd/subscriber/subscriber.go

probe:
	go build -o ./bin/probe/ ./cmd/probe/probe.go
	./bin/probe/probe

sub:
	go build -o ./bin/subscriber/ ./cmd/subscriber/subscriber.go
	./bin/subscriber/subscriber