.PHONY: dev test

dev:
	air

test:
	go test -v ./... -cover


swag:
	swag init -g cmd/homing/main.go


run:
	docker compose down
	docker image prune -f
	docker compose build app
	docker compose up --build -d
	docker compose ps
