# Stage 1: Build
FROM golang:1.23-alpine as builder

# Install dependencies
RUN apk add --no-cache git build-base

# Set the working directory inside the container
WORKDIR /app

# Copy dependency files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files and build
COPY . .
RUN go build -o main .

# Stage 2: Final Image
FROM alpine:latest

# Set up a non-root user
RUN adduser -D -u 1001 appuser

# Set the working directory inside the container
WORKDIR /app

# Copy only the necessary binary from the builder
COPY --from=builder /app/main .
COPY --from=builder /app/view /app/view
COPY --from=builder /app/static /app/static
COPY --from=builder /app/database /app/database
COPY --from=builder /app/controller /app/controller
COPY --from=builder /app/fileio /app/fileio
COPY --from=builder /app/helperfunc /app/helperfunc
COPY --from=builder /app/middlewares /app/middlewares
COPY --from=builder /app/pkg /app/pkg
COPY --from=builder /app/handlers /app/handlers
COPY --from=builder /app/model /app/model

# Set ownership and permissions for the non-root user
RUN chown -R appuser:appuser /app && chmod +x /app/main

# Switch to the non-root user
USER appuser

# Expose the application port
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["./main"]
