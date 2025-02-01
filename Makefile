build:
	@echo "Building application"
	@go build -o app
run:
	@echo "Starting application"
	@./app

dev: build run