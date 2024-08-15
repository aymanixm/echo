# Build the Go binary
FROM golang:alpine3.20 AS binary

# Install build dependencies
RUN apk add --no-cache build-base

# Set up your project
ADD . /app
WORKDIR /app

# Ensure dependencies are pulled in case of any
RUN go mod tidy

# Build the application
RUN go build -o http

# Create a lightweight image with the compiled binary
FROM alpine:3.20
WORKDIR /app
ENV PORT 8000
EXPOSE 8000
COPY --from=binary /app/http /app
CMD ["/app/http"]
