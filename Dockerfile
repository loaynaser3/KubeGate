# Use a lightweight Go image for building the application
FROM golang:1.23.4 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Build a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o kubegate main.go

# Use a compatible base image for running the application
FROM ubuntu:22.04

# Set the working directory
WORKDIR /app

# Install dependencies and kubectl
RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl && \
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && \
    rm kubectl && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Copy the built binary from the builder stage
COPY --from=builder /app/kubegate /usr/local/bin/kubegate


# Set the entrypoint
ENTRYPOINT ["kubegate"]

# Default command
CMD ["--help"]
