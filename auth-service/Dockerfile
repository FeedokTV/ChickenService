FROM golang:1.23.2 AS build

# Setting up working directory
WORKDIR /app

# Copy files into container
COPY . .

# Устанавливаем зависимости и строим приложение
RUN go mod tidy
# RUN go build -o auth-service ./cmd/main.go
RUN go build -o auth-service ./cmd/main.go

# Базовый контейнер для запуска приложения
#FROM gcr.io/distroless/base-
#COPY --from=build /app/auth-service /app/auth-service

# Порт, который будет открыт
EXPOSE 8081

# Запуск приложения
CMD ["/app/auth-service"]
