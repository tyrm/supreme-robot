fmt:
	go fmt ./...

lint:
	golint ./...

test-local: tidy fmt lint
	go test -cover -tags=postgres ./...

tidy:
	go mod tidy

.PHONY: fmt lint test