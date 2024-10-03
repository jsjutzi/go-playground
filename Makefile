setup:
	@echo "Checking Go installed..."
	@command -v go >/dev/null 2>&1 || { echo >&2 "Go is not installed.  Aborting."; exit 1; }

	@echo "Installing Go dependencies..."
	@go mod tidy

	@echo "Setting up pre-commit hooks..."
	@pre-commit install

	@echo "Setting up dev dependencies..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest


test:
	@go test ./...

build:
	@go build ./...

clean:
	@rm du-import-service

.PHONY:
.IGNORE:
