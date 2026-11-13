.PHONY: run mock

run:
	go run cmd/traceairlock/main.go

mock:
	go run cmd/mock-mcp/main.go