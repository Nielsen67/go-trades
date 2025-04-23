# Use the official Go image to build the app
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum (if they exist) and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .
COPY .env ./.env


# Build the Go app
RUN apk --no-cache add ca-certificates tzdata && \ 
    go build -o myapp .

# Use a smaller base image for the final container
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/myapp .
COPY --from=builder /app/.env .env


# Expose port 3000
EXPOSE 3000


# Command to run the app
CMD ["./myapp"]