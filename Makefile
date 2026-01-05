pull:
	@git pull origin master

build_and_compose:
	@docker compose up -d --build
compose:
	@docker compose up -d

cleanup:
	@echo "Cleaninig Docker cache"
	@docker system prune --all --force

update: pull build_and_compose cleanup

restart: compose


dev: format 
	@echo "Running on Windows"
	@go build -o app
	@./app
lint:
	@echo "Running linter"
	@golangci-lint run 
format:
	@echo "Fomtatting code"
	@golangci-lint fmt