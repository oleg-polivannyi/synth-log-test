# Start from the official Golang base image
FROM golang:1.23-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod file
COPY go.mod ./

# Download all dependencies. Dependencies will be cached if the go.mod file is not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Use a minimal base image for the final stage
FROM alpine:3.21 AS app

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the built Go app from the builder stage
COPY --from=builder /app/main ./main

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]