# The base go-image
FROM golang:1.15

# Set working directory
WORKDIR /app

# Copy all files from the current directory to the app directory
COPY . /app

# Build application
RUN go build -o bin/server server.go

# Expose port 5050
EXPOSE 5050

# Run the server executable
CMD ["bin/server"]