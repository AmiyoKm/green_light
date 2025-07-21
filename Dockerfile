FROM golang:1.24.5-alpine3.22

WORKDIR /usr/src/app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -ldflags='-s -w' -o=./bin/api ./cmd/api

# Expose port (adjust if your app uses a different port)
EXPOSE 8080

# Run the built binary instead of go run
CMD ["./bin/api"]
