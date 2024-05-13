# Use the official Golang image to build the binary.
# This ensures you have all necessary Go tools and libraries.
FROM golang:1.22.2 as builder
WORKDIR /app

# Copy everything from the current directory to the /app directory in the container
COPY ../app .

# Build the application. This results in a statically linked executable named 'app'.
# If needed, add flags like -ldflags="-s -w" to reduce binary size.
RUN CGO_ENABLED=0 GOOS=linux go build -o app

# Start from a smaller base image to run the app.
FROM alpine:latest  

# Add CA Certificates for HTTPS connections.
RUN apk --no-cache add ca-certificates

# Set the working directory to /root/
WORKDIR /root/

# Copy the built executable from the builder stage.
COPY --from=builder /app/app .

# Bind container to port 80.
EXPOSE 80

# Command to run the executable.
CMD ["./app"]
