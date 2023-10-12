# Stage 1: Build the Go application
FROM golang:1.20 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go mod and sum files to the container
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code to the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd/main

# Stage 2: Create a lightweight final image
FROM alpine:3.14

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/myapp .

# Run the application
CMD ["./myapp"]
