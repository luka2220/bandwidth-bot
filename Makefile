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
	$(info Running all of the application unit tests)
	@go test ./...

# Build the application binary
build:
	$(info Building the application binary into bin/)
	@go build -o bin/ ./cmd/main.go
