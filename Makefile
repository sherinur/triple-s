BINARY_NAME=triple-s

MAIN_FILE=cmd/main.go


build:
	gofumpt -l -w .
	go mod tidy
	go build -o $(BINARY_NAME) $(MAIN_FILE)

run: build
	./$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	gofumpt -l -w .
	go mod tidy
	rm -f $(BINARY_NAME)

help:
	@echo "Makefile commands:"
	@echo "  make build   - Build the project"
	@echo "  make run     - Build and run the project"
	@echo "  make clean   - Remove the compiled binary"
	@echo "  make help    - Show this help message"

commit:
	git add .
	git commit -m "Commit $$(date '+%Y-%m-%d %H:%M:%S')"
	git push

.PHONY: build run clean help
