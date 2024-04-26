SHELL := /bin/bash


# Test runs golang tests
test:
	@echo "Running tests..."
	go test .\...
	@echo "Finished tests..."

# Dev runs
dev:
	set MONGODB_URI=mongodb://127.0.0.1:27017&& set APP_PORT=8080&& go run .


tidy:
	go mod tidy