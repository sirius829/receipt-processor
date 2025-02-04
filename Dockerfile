# Build stage

FROM golang:1.22.11 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o receipt-processor ./cmd/receipt-processor

# Final Stage: Use a minimal base image.
FROM alpine:latest

# Create a non-root user.
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /home/appuser

# Copy the binary from the builder stage.
COPY --from=builder /app/receipt-processor .

# Change ownership to the non-root user.
RUN chown appuser:appgroup receipt-processor

# Expose the port.
EXPOSE 8080

# Switch to non-root user.
USER appuser

# Start the application.
CMD ["./receipt-processor"]