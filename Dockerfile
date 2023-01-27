

# Start from golang base image
FROM golang:alpine

# Add Maintainer info
LABEL maintainer="supachai phobpimai"
# wait 10s
# RUN sleep 10

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

# Setup folders
RUN mkdir /app
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY . ./
# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Build the Go app
RUN go build -o /docker_test

# Expose port 8080 to the outside world
EXPOSE 9000

# Run the executable
CMD [ "/docker_test" ]
