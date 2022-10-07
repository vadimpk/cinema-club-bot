FROM golang:1.19-alpine3.16 AS builder

COPY . /github.com/vadimpk/cinema-club-bot
WORKDIR /github.com/vadimpk/cinema-club-bot

RUN go mod download
RUN go build -o ./bin/bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/vadimpk/cinema-club-bo/bin/bot .
COPY --from=0 /github.com/vadimpk/cinema-club-bo/.env .
COPY --from=0 /github.com/vadimpk/cinema-club-bo/configs configs/

EXPOSE 80

CMD ["./bot"]