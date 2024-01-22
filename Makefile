# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Run the application
run:
	@templ generate 
	@go run cmd/api/main.go

templ:
	@templ generate 
	
bun:
	@bun build ./cmd/web/assets/js/app.ts --outdir ./cmd/web/assets/js
test:
	@echo "Testing..."
	@go test ./tests -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main
tw:
	@bunx tailwindcss -i ./cmd/web/css/style.css -o ./cmd/web/css/out.css --watch
# Live Reload
watch:
	@air
	# @if command -v air > /dev/null; then \
	#     air server --port 1337; \
	#     echo "Watching...";\
	# else \
	#     read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	#     if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	#         go install github.com/cosmtrek/air@latest; \
	#         air; \
	#         echo "Watching...";\
	#     else \
	#         echo "You chose not to install air. Exiting..."; \
	#         exit 1; \
	#     fi; \
	# fi

.PHONY: all build run test clean
