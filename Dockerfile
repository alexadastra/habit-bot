FROM golang:latest

# Set the current working directory inside the container
WORKDIR /app

# Copy the source code from the current directory to the working directory inside the container
COPY . .

# Install the dependencies
RUN go get -d -v ./...

# Build the app
RUN go build -o main .

# Run the app
CMD ["./main"]
