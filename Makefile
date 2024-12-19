# View the makefile commads
view:
	@cat Makefile

# Format code 
fmt:
	go fmt ./...

# View possible issues in codebase
vet:
	go vet ./...

# Add any missing libraries and remove unsed ones
tidy: fmt
	go mod tidy

# Run the serevr for testing
run:
	go run ./cmd

# Run all of the unit tests
test:
	@go test ./... -v

test-fwc:
	@go test fixed_window_test.go -v

test-tb:
	@go test token_bucket_test.go -v

# Build the application binary
build:
	@go build -o bin/ ./cmd/main.go
