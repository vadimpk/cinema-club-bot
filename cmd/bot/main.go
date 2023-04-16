package main

import "github.com/vadimpk/cinema-club-bot/internal/app"

const configsDir = "configs"
const configsFile = "local"

func main() {
	app.Run(configsDir, configsFile)
}
