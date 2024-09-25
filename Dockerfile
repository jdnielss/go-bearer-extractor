# Stage 1: Build stage
FROM golang:1.23-alpine AS buildstage

WORKDIR /app

# Copy the Go Modules and dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o go-bearer

# Stage 2: Final stage (smaller image with only the binary)
FROM alpine:3.14

WORKDIR /app

# Copy the compiled Go binary from the build stage
COPY --from=buildstage /app/go-bearer /app/go-bearer

# Set the binary as the entry point
ENTRYPOINT ["/app/go-bearer"]