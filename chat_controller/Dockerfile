# Use an official Golang runtime as a parent image
FROM golang:1.21

# Install git (required for go get)
RUN apt-get update && apt-get install -y git

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go app
RUN go build -o main .

# Command to run the executable
CMD ["./main"]
