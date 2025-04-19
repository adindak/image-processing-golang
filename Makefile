APP_NAME := image-processing-golang

.PHONY: build run clean test tidy

build:
	@echo "🔨 Building $(APP_NAME)..."
	go build -o $(APP_NAME) main.go

run:
	@echo "🚀 Running $(APP_NAME)..."
	go run main.go

clean:
	@echo "🧹 Cleaning up..."
	rm -f $(APP_NAME)

tidy:
	@echo "📦 Tidying up modules..."
	go mod tidy
