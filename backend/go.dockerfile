# Use the official Golang image as a base
FROM golang:1.22

# Set the working directory in the Docker container
WORKDIR /app

# Copy the Go files from your host to your Docker container
COPY . .

RUN apt-get update && apt-get install -y sqlite3

# Download and install the dependencies
RUN go get -d -v ./...

# Build the Go app
RUN go build -o /usr/local/bin/api ./cmd/main.go

# Make the 'api' file executable
# RUN chmod +x /usr/local/bin/api

# Expose port 8000 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/usr/local/bin/api"]
