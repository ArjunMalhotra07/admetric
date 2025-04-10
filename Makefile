run:
	@go run ./cmd/main.go

compose:
	@docker compose up --build

remove:
	@docker stop $$(docker ps -q) 2>/dev/null || true
	@docker rm -f $$(docker ps -aq) 2>/dev/null || true

stack:
	@docker compose -f stack-compose.yml up --build

compose:
	@docker compose up --build