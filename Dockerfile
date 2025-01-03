# Используем базовый образ Go
FROM golang:1.23-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Загружаем зависимости
COPY go.mod go.sum ./
RUN go mod tidy

# Сборка приложения
RUN go build -o main ./cmd/main.go

# Устанавливаем команду по умолчанию для запуска приложения
CMD ["./main"]