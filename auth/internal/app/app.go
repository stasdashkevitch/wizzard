package app

import (
	"os"
	"os/signal"
	"syscall"

	grpcapp "github.com/stasdashkevitch/wizzard/auth/internal/app/grpc"
	"github.com/stasdashkevitch/wizzard/common/logger"
)

type App struct {
	GRPCServer *grpcapp.App
	log        logger.ILogger
}

func New(
	log logger.ILogger,
	port int,
) *App {
	grpcApp := grpcapp.New(log, port)

	return &App{
		GRPCServer: grpcApp,
		log:        log,
	}
}

func (app *App) Run() {
	app.GRPCServer.MustRun()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	app.GRPCServer.Stop()
	app.log.Info("gracefully stopped")
}
