.PHONY: build run test clean docker-build docker-run docker-up docker-down test-coverage

# Build the application
build:
	go build -o bin/server ./cmd/server

# Run locally
run: build
	./bin/server

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -cover ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Build Docker image
docker-build:
	docker build -t order-pack-calculator .

# Run with Docker
docker-run: docker-build
	docker run -p 8080:8080 order-pack-calculator

# Run with docker-compose
docker-up:
	docker-compose up --build

# Stop docker-compose
docker-down:
	docker-compose down
