# Use the official Golang image as the build environment
FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY . .
# Build the Go application & docs
RUN go build -o application main.go

# Use a minimal Alpine-based image for the final image
FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
WORKDIR /build
# Copy the binary from the builder stage
COPY --from=builder /build/application /build/application
# Specify the command to run the application
RUN chmod +x application
ENTRYPOINT [ "/build/application"]