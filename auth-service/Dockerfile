FROM golang:1.23.2 AS build

# Setting up working directory
WORKDIR /app

# Copy files into container
COPY . .

RUN go mod tidy
RUN go build -o auth-service ./cmd/main.go

EXPOSE 8081

# Run application
CMD ["/app/auth-service"]
