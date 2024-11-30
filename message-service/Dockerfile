FROM golang:1.23.2 AS build

# Setting up working directory
WORKDIR /app

# Copy files into container
COPY . .

# Setup dependencies and making our app
RUN go mod tidy
RUN go build -o message-service ./cmd/main.go

# Basic container. Here im prefer not to use it :) But you cn
#FROM gcr.io/distroless/base-
#COPY --from=build /app/auth-service /app/auth-service

# Exposing port
EXPOSE 8081

# Run
CMD ["/app/message-service"]
