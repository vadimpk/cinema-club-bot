package main

import "github.com/vadimpk/cinema-club-bot/internal/app"

const configsDir = "configs"
const configsFile = "prod"

func main() {
	app.Run(configsDir, configsFile)
}
