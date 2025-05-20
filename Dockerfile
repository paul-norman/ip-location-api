# Use a multi-stage build for a smaller final image
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache make git

# Set up working directory
WORKDIR /app

# Copy only the files needed for dependency installation first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build for both architectures to support different platforms
RUN make build_linux_amd64 && make build_linux_arm64

# Create a minimal runtime image
FROM --platform=$TARGETPLATFORM alpine:3.19 AS runtime

# Install CA certificates for HTTPS requests
RUN apk add --no-cache ca-certificates tzdata

# Create a non-root user to run the application
RUN adduser -D -h /app appuser

# Create directories for application data
WORKDIR /app
RUN mkdir -p /app/downloads && \
    chown -R appuser:appuser /app

# Create the entrypoint script that will generate .env
RUN echo '#!/bin/sh' > /app/entrypoint.sh && \
    echo 'echo "SERVER_HOST=$SERVER_HOST" > /app/.env' >> /app/entrypoint.sh && \
    echo 'echo "SERVER_PORT=$SERVER_PORT" >> /app/.env' >> /app/entrypoint.sh && \
    echo 'echo "API_KEY=$API_KEY" >> /app/.env' >> /app/entrypoint.sh && \
    echo 'echo "COUNTRY=$COUNTRY" >> /app/.env' >> /app/entrypoint.sh && \
    echo 'echo "CITY=$CITY" >> /app/.env' >> /app/entrypoint.sh && \
    echo 'echo "ASN=$ASN" >> /app/.env' >> /app/entrypoint.sh && \
    echo 'echo "UPDATE_TIME=$UPDATE_TIME" >> /app/.env' >> /app/entrypoint.sh && \
    echo 'echo "LOAD_LOG_FREQ=$LOAD_LOG_FREQ" >> /app/.env' >> /app/entrypoint.sh && \
    echo 'echo "DB_TYPE=$DB_TYPE" >> /app/.env' >> /app/entrypoint.sh && \
    echo '' >> /app/entrypoint.sh && \
    echo '# Add database variables if they are set' >> /app/entrypoint.sh && \
    echo 'if [ ! -z "$DB_HOST" ]; then' >> /app/entrypoint.sh && \
    echo '  echo "DB_HOST=$DB_HOST" >> /app/.env' >> /app/entrypoint.sh && \
    echo 'fi' >> /app/entrypoint.sh && \
    echo 'if [ ! -z "$DB_PORT" ]; then' >> /app/entrypoint.sh && \
    echo '  echo "DB_PORT=$DB_PORT" >> /app/.env' >> /app/entrypoint.sh && \
    echo 'fi' >> /app/entrypoint.sh && \
    echo 'if [ ! -z "$DB_USER" ]; then' >> /app/entrypoint.sh && \
    echo '  echo "DB_USER=$DB_USER" >> /app/.env' >> /app/entrypoint.sh && \
    echo 'fi' >> /app/entrypoint.sh && \
    echo 'if [ ! -z "$DB_PASS" ]; then' >> /app/entrypoint.sh && \
    echo '  echo "DB_PASS=$DB_PASS" >> /app/.env' >> /app/entrypoint.sh && \
    echo 'fi' >> /app/entrypoint.sh && \
    echo 'if [ ! -z "$DB_NAME" ]; then' >> /app/entrypoint.sh && \
    echo '  echo "DB_NAME=$DB_NAME" >> /app/.env' >> /app/entrypoint.sh && \
    echo 'fi' >> /app/entrypoint.sh && \
    echo 'if [ ! -z "$DB_SCHEMA" ]; then' >> /app/entrypoint.sh && \
    echo '  echo "DB_SCHEMA=$DB_SCHEMA" >> /app/.env' >> /app/entrypoint.sh && \
    echo 'fi' >> /app/entrypoint.sh && \
    echo '' >> /app/entrypoint.sh && \
    echo '# Start the application' >> /app/entrypoint.sh && \
    echo 'exec /app/ip-location-api' >> /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh && \
    chown appuser:appuser /app/entrypoint.sh

USER appuser

# Create separate stages for each architecture
FROM runtime AS amd64
COPY --from=builder /app/builds/ip-location-api-linux-x64.bin /app/ip-location-api

FROM runtime AS arm64
COPY --from=builder /app/builds/ip-location-api-linux-arm64.bin /app/ip-location-api

# Use the appropriate image based on architecture
FROM ${TARGETARCH}

# Declare a volume for persisting downloaded data
VOLUME /app/downloads

# Expose the server port (will be overridden by SERVER_PORT env var if set)
EXPOSE 8080

# Set up environment variables with defaults
# These can be overridden when running the container
ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=8080
ENV API_KEY=""
ENV COUNTRY="geo-whois-asn-country"
ENV CITY=""
ENV ASN="asn"
ENV UPDATE_TIME="01:30"
ENV LOAD_LOG_FREQ=50000
# DB_TYPE can be mmdb, postgres, mysql, sqlite or :memory:
ENV DB_TYPE="mmdb"
# Database connection variables (used when DB_TYPE is set to postgres/mysql/sqlite)
# ENV DB_HOST=""
# ENV DB_PORT=""
# ENV DB_USER="" # used with sqlite for .db filename
# ENV DB_PASS=""
# ENV DB_NAME=""
# ENV DB_SCHEMA="" # used for postgres/sqlite


# Command to run the entrypoint script that will generate .env and start the app
CMD ["/app/entrypoint.sh"]
