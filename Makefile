SHELL := /bin/bash


# Test runs golang tests
test:
	@echo "Running tests..."
	go test .\...
	@echo "Finished tests..."

# Dev runs
dev:
	export MONGODB_URI=mongodb://openbank:openbank@127.0.0.1:27017&& export APP_PORT=8080&& go run .


tidy:
	go mod tidy