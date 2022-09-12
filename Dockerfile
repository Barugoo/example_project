# Start from golang base image
FROM golang:alpine

# Setup folders
RUN mkdir /app
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o /build

# Expose port 8080 and 8081 to the outside world
EXPOSE 8080
EXPOSE 8081

# Run the executable
CMD [ "/build" ]