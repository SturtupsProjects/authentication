package main

import (
	"authentification/config"
	"authentification/internal/app"
)

func main() {
	cfg := config.NewConfig()

	app.Run(&cfg)
}
