.PHONY: dev test

dev:
	air

test:
	go test -v ./... -cover