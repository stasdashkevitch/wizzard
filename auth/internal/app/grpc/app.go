package grpcapp

import (
	"fmt"
	"net"

	authgrpc "github.com/stasdashkevitch/wizzard/auth/internal/controller/grpc/auth"
	"github.com/stasdashkevitch/wizzard/common/logger"
	"google.golang.org/grpc"
)

type App struct {
	log        logger.ILogger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log logger.ILogger,
	port int,
) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *App) Run() error {
	const op = "grpcapp.Run"

	app.log.Info("starting grpc server")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	app.log.Info("grpc server is running", "addr", l.Addr().String())

	if err := app.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (app *App) Stop() {
	const op = "grpcapp.Stop"

	app.log.Info("stopping GRPC server", "port", app.port)

	app.gRPCServer.Stop()
}
