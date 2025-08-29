
test:
	@echo "Running tests..."
	go test ./...

lint:
	@echo "Running golangci-lint locally..."
	golangci-lint run

lint-fix:
	@echo "Running golangci-lint with fix..."
	golangci-lint run --timeout=5m --fix

help:
	@echo "Makefile targets:"
	@echo "  test    			- Run tests"
	@echo "  lint         		- Run golangci-lint locally"
	@echo "  lint-fix     		- Run golangci-lint and automatically fix issues"
	@echo "  help    			- Show this help message"