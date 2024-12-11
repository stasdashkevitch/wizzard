package main

import (
	"github.com/stasdashkevitch/wizzard/auth/internal/app"
	"github.com/stasdashkevitch/wizzard/auth/internal/config"
	"github.com/stasdashkevitch/wizzard/common/logger"
)

func main() {
	cfg := config.NewConfig()
	log := logger.NewLogger(cfg.Env)

	log.Info("starting app")

	app := app.New(log, cfg.Port)
	app.Run()
}
