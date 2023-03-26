.PHONY: test run

test:
	go test ./tests/...

run:
	go run cmd/main.go
