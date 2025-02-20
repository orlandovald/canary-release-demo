FROM --platform=$BUILDPLATFORM golang:1.24.0 AS builder

# Set up environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /app

# Copy source files
COPY . .

# Detect target architecture and set GOARCH accordingly
ARG TARGETARCH
ENV GOARCH=$TARGETARCH

# Build the binary with optimization flags
RUN go build -ldflags "-s -w" -o server

# Final image
FROM --platform=$TARGETPLATFORM scratch

# Copy the compiled binary
COPY --from=builder /app/server /server

# Copy SSL certificates for HTTPS support
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

CMD ["/server"]

