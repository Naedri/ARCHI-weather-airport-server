build:
	go build -o ./bin/probe/ ./cmd/probe/probe.go
	go build -o ./bin/subscriber/ ./cmd/subscriber/subscriber.go

probe:
	go build -o ./bin/probe/ ./cmd/probe/probe.go
	cd ./bin/probe/ && \
	./probe

sub:
	go build -o ./bin/subscriber/ ./cmd/subscriber/subscriber.go
	cd ./bin/subscriber/ && \
	./subscriber

http:
	go build -o ./bin/httpServer ./cmd/httpServer/main.go
	cd ./bin/httpServer/ && \
	./httpServer