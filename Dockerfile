# Use a full-featured base image for building
FROM golang:1.21.6-alpine3.18 AS builder

# Install git if your project requires
RUN apk update && apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /src

# Copy the source code into the container
COPY . .

# Fetch dependencies using go mod if your project uses Go modules
RUN go mod download

# Version and Git Commit build arguments
ARG VERSION="0.0.0"
ARG GIT_COMMIT="unknown"
ARG BUILD_DATE=""

# Create a non-root user and group 'cache' (use UID and GID to be compatible with scratch)
# Note: Alpine uses 'addgroup' and 'adduser' instead of 'groupadd' and 'useradd'
RUN addgroup -S cache && adduser -S cache -G cache

# Build the Go app with versioning information
RUN GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/supporttools/go-web-cache/pkg/health.version=$VERSION -X github.com/supporttools/go-web-cache/pkg/health.GitCommit=$GIT_COMMIT -X github.com/supporttools/go-web-cache/pkg/health.BuildTime=$BUILD_DATE" -o /bin/go-web-cache

# Start from scratch for the runtime stage
FROM scratch

# Set the working directory to /app
WORKDIR /app

# Copy the built binary and config file from the builder stage
COPY --from=builder /bin/go-web-cache /app/
COPY config.json /app/

# Import the user and group files from the builder stage
# This step is necessary to use the non-root user created in the builder stage
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Use the non-root user 'cache'
USER cache

# Set the binary as the entrypoint of the container
ENTRYPOINT ["/app/go-web-cache"]
