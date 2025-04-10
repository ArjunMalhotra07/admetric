run:
	@go run ./cmd/main.go

remove:
	@docker stop $$(docker ps -q) 2>/dev/null || true
	@docker rm -f $$(docker ps -aq) 2>/dev/null || true

compose:
	@docker compose up --build

client:
	@go run ./client/click_simulator.go