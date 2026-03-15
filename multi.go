package slogutil

import (
	"context"
	"log/slog"
)

type multiHandler struct {
	handlers []slog.Handler
}

// Multi returns a slog.Handler that fans out log records to all provided handlers.
// Enabled returns true if any underlying handler is enabled for the given level.
// Handle sends each record to every handler; the first error encountered is returned.
func Multi(handlers ...slog.Handler) slog.Handler {
	all := make([]slog.Handler, len(handlers))
	copy(all, handlers)
	return &multiHandler{handlers: all}
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}
