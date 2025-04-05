.PHONY: run
run:
	@go run .

.PHONY: up
up:
	@docker compose up 

.PHONY: tidy
tidy:
	@go mod tidy