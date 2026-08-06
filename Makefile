default: test

.PHONY: test test-full generate bench benchmark release-dryrun

test:
	@echo "Running short test suite..."
	@go test -short -tags magicgen ./...

test-full:
	@echo "Running full test suite..."
	@go test -tags magicgen ./...

release-dryrun:
	@echo "Running release dry-run..."
	goreleaser release --snapshot --clean

generate:
	@echo "Generating code..."
	@go generate ./...

bench:
	@echo "Running benc..."
	go run . bench

benchmark:
	@echo "Running benchmarks..."
	go test -run=^$$ -bench=. -benchmem ./...
