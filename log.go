package amplifyx

import (
	"context"
	"io"
	"log/slog"
)

type logHandler struct {
	normal   slog.Handler
	abnormal slog.Handler
}

func newLogHandler(normal, abnormal slog.Handler) *logHandler {
	return &logHandler{
		normal:   normal,
		abnormal: abnormal,
	}
}

func (h logHandler) handler(level slog.Level) slog.Handler {
	if level >= slog.LevelWarn {
		return h.abnormal
	}

	return h.normal
}

func (h logHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler(level).Enabled(ctx, level)
}

func (h logHandler) Handle(ctx context.Context, record slog.Record) error {
	return h.handler(record.Level).Handle(ctx, record)
}

func (h logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return newLogHandler(
		h.normal.WithAttrs(attrs),
		h.abnormal.WithAttrs(attrs),
	)
}

func (h logHandler) WithGroup(name string) slog.Handler {
	return newLogHandler(
		h.normal.WithGroup(name),
		h.abnormal.WithGroup(name),
	)
}

func NewLogger(normalOut, abnormalOut io.Writer) *slog.Logger {
	return slog.New(newLogHandler(
		slog.NewJSONHandler(normalOut, nil),
		slog.NewJSONHandler(abnormalOut, nil),
	))
}
