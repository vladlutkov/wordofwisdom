FROM golang:1.24-alpine AS builder

WORKDIR /bin

COPY ./src .

RUN go mod download

RUN go build -o server .

# Используем минимальный образ для продакшн-контейнера
FROM alpine:latest

WORKDIR /root/

# Копируем скомпилированный бинарник
COPY --from=builder /bin/server .
# и нашу БД
COPY --from=builder /bin/quotes.txt .

# Запускаем сервис
ENV PORT=44542
CMD ["./server"]
