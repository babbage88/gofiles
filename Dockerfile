# Build stage
FROM golang:1.23.4 AS builder

WORKDIR /app

# Copy and download dependency using go mod
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o fileserver .

# Final stage copy bin and install pre-requisites
FROM alpine:latest AS final

WORKDIR /app

# Copy the built binary from the builder
COPY --from=builder /app/fileserver /app/fileserver

# Copy any additional files (e.g., templates, CSS, etc.)
COPY templates/ /app/templates/
COPY main.css /app/

# Expose the port your app runs on
EXPOSE 8080

# Command to run the executable
CMD ["/app/fileserver"]
