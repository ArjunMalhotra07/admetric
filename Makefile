.PHONY: run
run:
	@go run ./cmd/main.go

.PHONY: up
up:
	@docker compose down
	@docker compose up --build

.PHONY: remove
remove:
	@echo "Stopping all running containers..."
	@docker stop $$(docker ps -q) 2>/dev/null || true
	@echo "Removing all containers..."
	@docker rm -f $$(docker ps -aq) 2>/dev/null || true
	@echo "Done."

.PHONY: tidy
tidy:
	@go mod tidy