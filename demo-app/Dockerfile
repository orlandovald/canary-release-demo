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

ARG ERROR_RATE=0
ARG APP_VERSION=dev

# Build the binary with optimization flags
RUN go build -ldflags="-s -w -X 'main.version=${APP_VERSION}' -X 'main.errorRate=${ERROR_RATE}'" -o server

# Final image
FROM --platform=$TARGETPLATFORM scratch

# Copy the compiled binary
COPY --from=builder app/server /server
COPY static /static

EXPOSE 8080

CMD ["/server"]

