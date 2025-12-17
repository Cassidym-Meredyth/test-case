# Этап сборки веб-приложения
FROM golang:1.25-alpine3.23 AS builder

# Рабочая директория
WORKDIR /app

# Копирование модулей
COPY go.mod go.sum ./
RUN go mod download

# Копирование основного кода + html  
COPY ./html/* main.go ./

# Сборка статически линкованного бинарника
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o weather-app main.go

# Этап запуска веб-приложения
FROM alpine:3.23

WORKDIR /app

# Копирование бинарника, шаблона и css стилей из builder
COPY --from=builder /app/weather-app .
COPY --from=builder /app/index.html /app/html/
COPY --from=builder /app/style.css /app/html/style/

# Порт прослушивания
EXPOSE 8080

CMD ["./weather-app"]