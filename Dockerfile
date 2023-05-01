FROM golang:alpine

# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

# Setup folders
RUN mkdir /shopping-app
WORKDIR /shopping-app

# Copy the source from the current directory to the working Directory inside the container
COPY . .
COPY .env .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
ENTRYPOINT ["go","run","./cmd/api/","."]