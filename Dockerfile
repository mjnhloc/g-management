# Use the official Golang image
FROM golang:1.23

# Set the working directory
WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main ./cmd/main.go

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["/app/main"]
