.PHONY:
.SILENT:

build:
	go build -o ./.bin/bot cmd/bot/main.go

start-redis:
	docker run --name redis-test -p 6379:6379 -d redis --requirepass "redis-test"

run: build
	./.bin/bot
