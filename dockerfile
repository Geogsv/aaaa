# Используем официальный образ Golang
FROM golang:1.21 as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем Go модули и зависимости
COPY go.mod ./
RUN go mod download

# Копируем всё остальное
COPY . .

# Собираем бинарник
RUN go build -o main .

# Финальный образ
FROM debian:bookworm-slim

# Создаем директории
WORKDIR /app

# Копируем бинарник и статику из предыдущего этапа
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
COPY --from=builder /app/contacts.json ./contacts.json

# Открываем порт (если нужно)
EXPOSE 8080

# Запускаем сервер
CMD ["./main"]