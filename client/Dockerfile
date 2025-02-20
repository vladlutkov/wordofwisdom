FROM golang:1.24-alpine AS builder

WORKDIR /bin

COPY ./src .

RUN go mod download

RUN go build -o client .

# Финальный этап: минимальный образ
FROM alpine:latest

WORKDIR /root/

# Копируем скомпилированный бинарник
COPY --from=builder /bin/client .

# Запускаем сервис
ENV ADDR="127.0.0.1:44542"
CMD ["./client"]
