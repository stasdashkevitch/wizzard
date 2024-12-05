package logger

import (
	"log/slog"
	"os"
)

type ILogger interface {
	Debug()
}

type Logger struct {
	logger *slog.Logger
}

const (
	LevelTrace   = slog.Level(-8)
	LevelDebug   = slog.LevelDebug
	LevelInfo    = slog.LevelInfo
	LevelWarning = slog.LevelWarn
	LevelError   = slog.LevelError
	LevelFatal   = slog.Level(12)
)

var LevelNames = map[slog.Leveler]string{
	LevelTrace:   "TRACE",
	LevelDebug:   "DEBUG",
	LevelInfo:    "INFO",
	LevelWarning: "WARNING",
	LevelError:   "ERROR",
	LevelFatal:   "FATAL",
}

func New(level string) *Logger {
	var l slog.Level

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     l,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]
				if !exists {
					levelLabel = level.String()
				}

				a.Value = slog.StringValue(levelLabel)
			}

			return a
		},
	}

	prettyHandlerOptions := PrettyHandlerOptions{
		SlogOpts: *opts,
	}
	handler := NewPrettyHandler(os.Stdout, prettyHandlerOptions)
	logger := slog.New(handler)

	return &Logger{
		logger: logger,
	}
}
