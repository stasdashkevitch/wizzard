package main

import (
	"log/slog"

	"github.com/stasdashkevitch/wizzard/pkg/logger"
)

func main() {
	l := logger.NewLogger("Trace")
	l.Info("Hello")
	l.Info("Hello", slog.String("key", "value"))
	l.Debug("E")
	l.Trace("yee")
	l.Fatal("fatal")
}
