BINARY_NAME=triple-s

MAIN_FILE=cmd/main.go


build:
	@echo "Building the project..."
	go build -o $(BINARY_NAME) $(MAIN_FILE)

run: build
	@echo "Running the project..."
	./$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)

help:
	@echo "Makefile commands:"
	@echo "  make build   - Build the project"
	@echo "  make run     - Build and run the project"
	@echo "  make clean   - Remove the compiled binary"
	@echo "  make help    - Show this help message"

.PHONY: build run clean help
