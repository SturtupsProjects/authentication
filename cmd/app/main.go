package main

import (
	"authentification/internal/app"
	"authentification/config"
)

func main() {
	cfg := config.NewConfig()

	app.Run(cfg)
}
