.PHONY:
.SILENT:

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

docker-build:
	docker build -t bot .

docker-run: docker-build
	docker run -p 9090:80 --env-file=.env --name cinema-club-bot bot