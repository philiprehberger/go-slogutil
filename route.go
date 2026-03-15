package slogutil

import (
	"context"
	"log/slog"
)

type levelRouteHandler struct {
	routes map[slog.Level]slog.Handler
}

// LevelRoute returns a slog.Handler that routes records to different handlers
// based on the log level. If an exact level match is not found in routes, it
// falls back to the next lower level's handler. If no handler matches, the
// record is silently dropped.
//
// Common usage routes errors to a stderr handler and info messages to stdout:
//
//	slog.New(slogutil.LevelRoute(map[slog.Level]slog.Handler{
//	    slog.LevelInfo:  stdoutHandler,
//	    slog.LevelError: stderrHandler,
//	}))
func LevelRoute(routes map[slog.Level]slog.Handler) slog.Handler {
	cp := make(map[slog.Level]slog.Handler, len(routes))
	for k, v := range routes {
		cp[k] = v
	}
	return &levelRouteHandler{routes: cp}
}

func (lr *levelRouteHandler) Enabled(ctx context.Context, level slog.Level) bool {
	h := lr.resolve(level)
	if h == nil {
		return false
	}
	return h.Enabled(ctx, level)
}

func (lr *levelRouteHandler) Handle(ctx context.Context, r slog.Record) error {
	h := lr.resolve(r.Level)
	if h == nil {
		return nil
	}
	return h.Handle(ctx, r)
}

func (lr *levelRouteHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	routes := make(map[slog.Level]slog.Handler, len(lr.routes))
	for k, v := range lr.routes {
		routes[k] = v.WithAttrs(attrs)
	}
	return &levelRouteHandler{routes: routes}
}

func (lr *levelRouteHandler) WithGroup(name string) slog.Handler {
	routes := make(map[slog.Level]slog.Handler, len(lr.routes))
	for k, v := range lr.routes {
		routes[k] = v.WithGroup(name)
	}
	return &levelRouteHandler{routes: routes}
}

// resolve finds the handler for the given level, falling back to the next
// lower defined level.
func (lr *levelRouteHandler) resolve(level slog.Level) slog.Handler {
	// Try exact match first.
	if h, ok := lr.routes[level]; ok {
		return h
	}
	// Fall back to the highest level that is still <= the given level.
	var best slog.Handler
	bestLevel := slog.Level(-128)
	for l, h := range lr.routes {
		if l <= level && l > bestLevel {
			best = h
			bestLevel = l
		}
	}
	return best
}
