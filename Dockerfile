# Используем базовый образ Go
FROM golang:1.23.2-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Загружаем зависимости
RUN go mod tidy

# Сборка приложения
RUN go build -o main .

# Устанавливаем команду по умолчанию
CMD ["./main"]
