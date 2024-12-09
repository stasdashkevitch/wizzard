package logger

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/mdobak/go-xerrors"
)

type ILogger interface {
	Trace(msg string, args ...any)
	TraceContext(ctx context.Context, msg string, args ...any)
	Debug(msg string, args ...any)
	DebugContext(ctx context.Context, msg string, args ...any)
	Info(msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	Warn(msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	Error(msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
	Fatal(msg string, args ...any)
	FatalContext(ctx context.Context, msg string, args ...any)
}

type Logger struct {
	l *slog.Logger
}

var _ ILogger = (*Logger)(nil)

const (
	LevelTrace = slog.Level(-8)
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
	LevelFatal = slog.Level(12)
)

var LevelNames = map[slog.Leveler]string{
	LevelTrace: "TRACE",
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
}

type stackFrame struct {
	Func   string `json:"func"`
	Source string `json:"source"`
	Line   int    `json:"line"`
}

func marshalStack(err error) []stackFrame {
	trace := xerrors.StackTrace(err)

	if len(trace) == 0 {
		return nil
	}

	frames := trace.Frames()

	s := make([]stackFrame, len(frames))

	for i, v := range frames {
		f := stackFrame{
			Source: filepath.Join(
				filepath.Base(filepath.Dir(v.File)),
				filepath.Base(v.File),
			),
			Func: filepath.Base(v.Function),
			Line: v.Line,
		}

		s[i] = f
	}

	return s
}

func fmtErr(err error) slog.Value {
	var groupValues []slog.Attr

	groupValues = append(groupValues, slog.String("msg", err.Error()))

	frames := marshalStack(err)

	if frames != nil {
		groupValues = append(groupValues, slog.Any("trace", frames))
	}

	return slog.GroupValue(groupValues...)
}

func NewLogger(level string) ILogger {
	var l slog.Level

	switch strings.ToLower(level) {
	case "error":
		l = LevelError
	case "warn":
		l = LevelWarn
	case "info":
		l = LevelInfo
	case "debug":
		l = LevelDebug
	case "trace":
		l = LevelTrace
	}

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
			} else if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05"))
			}

			switch a.Value.Kind() {
			case slog.KindAny:
				switch v := a.Value.Any().(type) {
				case error:
					a.Value = fmtErr(v)
				}
			}

			return a
		},
	}

	handler := slog.Handler(slog.NewJSONHandler(os.Stdout, opts))
	logger := slog.New(handler)

	return &Logger{
		l: logger,
	}
}

func (l *Logger) Trace(msg string, args ...any) {
	err := xerrors.New(msg)
	l.l.Log(context.Background(), LevelTrace, msg, slog.Any("error", err))

}

func (l *Logger) TraceContext(ctx context.Context, msg string, args ...any) {
	err := xerrors.New(msg)
	l.l.Log(ctx, LevelTrace, msg, slog.Any("error", err))
}

func (l *Logger) Debug(msg string, args ...any) {
	l.l.Debug(msg, args...)
}

func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.l.DebugContext(ctx, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.l.Info(msg, args...)
}

func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.l.InfoContext(ctx, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.l.Warn(msg, args...)
}

func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.l.WarnContext(ctx, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.l.Error(msg, args...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.l.ErrorContext(ctx, msg, args...)
}

func (l *Logger) Fatal(msg string, args ...any) {
	l.l.Log(context.Background(), LevelFatal, msg, args...)
	os.Exit(1)
}

func (l *Logger) FatalContext(ctx context.Context, msg string, args ...any) {
	l.l.Log(ctx, LevelFatal, msg, args...)
	os.Exit(1)
}
