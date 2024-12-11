package main

import (
	"github.com/stasdashkevitch/wizzard/auth/internal/config"
	"github.com/stasdashkevitch/wizzard/common/logger"
)

func main() {
	cfg := config.NewConfig()
	logger := logger.NewLogger(cfg.Env)

	logger.Info("start app")
}
