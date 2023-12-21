# Stage 1: Build the Go binary
FROM golang:1.21.5 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

# Set environment variables for the build
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go mod download

COPY . .

RUN go build -o app

# Stage 2: Create the final lightweight image
FROM alpine:latest

WORKDIR /app

# Create a non-root user
RUN adduser -D -g '' appuser

COPY --from=builder /app/app .
RUN chown appuser:appuser app

USER appuser

EXPOSE 8080

CMD ["./app"]
