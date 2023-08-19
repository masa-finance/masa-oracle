# Use the official Golang image as our base image
FROM golang:1.20

# Install git (required for fetching dependencies)
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

# Fetch dependencies
RUN go mod tidy

# Build the Go app
RUN go build -o main .

# Expose port 4001 (change if necessary)
EXPOSE 4001

# Command to run the executable
CMD ["./main"]
