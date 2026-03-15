package slogutil

import (
	"context"
	"log/slog"
	"sync/atomic"
)

type samplingHandler struct {
	handler slog.Handler
	rate    int64
	counter atomic.Int64
}

// Sampling returns a slog.Handler that only passes 1 in every rate records to
// the underlying handler. This is useful for high-volume code paths where logging
// every event would be too expensive. Records at slog.LevelError or above are
// always logged regardless of the sampling rate. The counter is incremented
// atomically and is safe for concurrent use.
func Sampling(handler slog.Handler, rate int) slog.Handler {
	if rate <= 1 {
		return handler
	}
	return &samplingHandler{
		handler: handler,
		rate:    int64(rate),
	}
}

func (s *samplingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return s.handler.Enabled(ctx, level)
}

func (s *samplingHandler) Handle(ctx context.Context, r slog.Record) error {
	// Always log errors.
	if r.Level >= slog.LevelError {
		return s.handler.Handle(ctx, r)
	}
	n := s.counter.Add(1)
	if n%s.rate == 0 {
		return s.handler.Handle(ctx, r)
	}
	return nil
}

func (s *samplingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &samplingHandler{
		handler: s.handler.WithAttrs(attrs),
		rate:    s.rate,
	}
}

func (s *samplingHandler) WithGroup(name string) slog.Handler {
	return &samplingHandler{
		handler: s.handler.WithGroup(name),
		rate:    s.rate,
	}
}
