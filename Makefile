default: test

.PHONY: test generate

test:
	@echo "Running tests..."
	@go test -tags magicgen ./...

generate:
	@echo "Generating code..."
	@go generate ./...
