default: test

.PHONY: test test-full generate bench

test:
	@echo "Running short test suite..."
	@go test -short -tags magicgen ./...

test-full:
	@echo "Running full test suite..."
	@go test -tags magicgen ./...

generate:
	@echo "Generating code..."
	@go generate ./...

bench:
	@echo "Running benchmarks..."
	go test -run=^$$ -bench=. -benchmem ./...
