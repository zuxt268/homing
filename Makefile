.PHONY: dev test

dev:
	air

test:
	go test -v ./... -cover


swag:
	swag init -g cmd/homing/main.go
