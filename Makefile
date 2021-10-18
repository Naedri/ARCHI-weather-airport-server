build:
	go build -o ./bin/probe/ ./cmd/probe/probe.go
	go build -o ./bin/subscriber/ ./cmd/subscriber/subscriber.go

probe:
	go build -o ./bin/probe/ ./cmd/probe/probe.go
	cd ./bin/probe/ && \
	./probe

probe2:
	go build -o ./bin/probe2/ ./cmd/probe/probe.go
	cd ./bin/probe2/ && \
	./probe

sub:
	go build -o ./bin/subscriber/ ./cmd/subscriber/subscriber.go
	cd ./bin/subscriber/ && \
	./subscriber

exporter:
	go build -o ./bin/exporter/ ./cmd/exporter/exporter.go
	cd ./bin/exporter/ && \
	./exporter

http:
	go build -o ./bin/httpServer/ ./cmd/httpServer/httpServer.go
	cd ./bin/httpServer/ && \
	./httpServer