# Stage 1: Build the Go application
FROM golang:1.20-alpine AS builder


WORKDIR /app

ADD . /app
# Download and cache Go dependencies
RUN go mod tidy

RUN go build -o ticket_bin cmd/main.go

# Stage 2: Create a lightweight image to run the application
FROM alpine:3.14

# Set the working directory inside the container
WORKDIR /app

# Copy only the built binary from the previous stage
COPY --from=builder /app/ticket_bin .
COPY --from=builder /app/config/config.yaml /app/config/config.yaml
COPY --from=builder /app/internal/data/migrations /app/internal/data/migrations
COPY --from=builder /app/public /app/public
COPY --from=builder /app/makefile /app/makefile
# Expose any necessary ports
EXPOSE 8080

# Command to run the executable
CMD ["./ticket_bin"]
