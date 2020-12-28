FROM golang:1.15.6-alpine3.12

WORKDIR /build

# Copy and download dependencies (not really important right now)
COPY go.mod .
RUN go mod download

# Copy the source for the hashServer
COPY cmd/hashServer .

# Build the application
RUN go build -o main

# Copy the program to the dist folder
WORKDIR /dist
RUN cp /build/main .

#Expose the listen port
EXPOSE 8080

# Run the app when the container starts
CMD ["/dist/main"]