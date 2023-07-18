# Use the official Golang image from the Docker Hub
FROM golang:1.17

# Create a directory for the application inside the Docker image
WORKDIR /app

# Copy everything from the current directory on your machine to the app directory in the Docker image
COPY . .

# Download necessary dependencies
RUN go mod download

# This command runs your application when the Docker container is launched.
CMD ["go", "run", "mid.go"]
