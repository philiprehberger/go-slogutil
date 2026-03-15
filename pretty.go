// Package slogutil provides handlers, formatters, and helpers for log/slog.
package slogutil

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
)

// PrettyHandlerOptions configures the PrettyHandler.
type PrettyHandlerOptions struct {
	// Level is the minimum level to log. If nil, defaults to slog.LevelInfo.
	Level slog.Leveler

	// TimeFormat is the time layout string. If empty, defaults to time.DateTime.
	TimeFormat string
}

type prettyHandler struct {
	opts   PrettyHandlerOptions
	w      io.Writer
	mu     *sync.Mutex
	attrs  []slog.Attr
	groups []string
}

// PrettyHandler returns a new slog.Handler that writes colorized, human-readable
// log lines to w. It is intended for development console output.
func PrettyHandler(w io.Writer, opts *PrettyHandlerOptions) slog.Handler {
	h := &prettyHandler{
		w:  w,
		mu: &sync.Mutex{},
	}
	if opts != nil {
		h.opts = *opts
	}
	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}
	if h.opts.TimeFormat == "" {
		h.opts.TimeFormat = time.DateTime
	}
	return h
}

func (h *prettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *prettyHandler) Handle(_ context.Context, r slog.Record) error {
	timeStr := r.Time.Format(h.opts.TimeFormat)
	levelColor := levelToColor(r.Level)
	levelStr := r.Level.String()

	h.mu.Lock()
	defer h.mu.Unlock()

	// Write timestamp, level, and message.
	_, err := fmt.Fprintf(h.w, "%s %s%s%s %s",
		timeStr, levelColor, levelStr, colorReset, r.Message)
	if err != nil {
		return err
	}

	// Write pre-set attrs.
	for _, a := range h.attrs {
		if err := h.writeAttr(a); err != nil {
			return err
		}
	}

	// Write record attrs.
	r.Attrs(func(a slog.Attr) bool {
		err = h.writeAttr(a)
		return err == nil
	})
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(h.w)
	return err
}

func (h *prettyHandler) writeAttr(a slog.Attr) error {
	if a.Equal(slog.Attr{}) {
		return nil
	}
	prefix := ""
	for _, g := range h.groups {
		prefix += g + "."
	}
	if a.Value.Kind() == slog.KindGroup {
		for _, ga := range a.Value.Group() {
			nested := slog.Attr{Key: prefix + a.Key + "." + ga.Key, Value: ga.Value}
			if err := h.writeAttrFlat(nested); err != nil {
				return err
			}
		}
		return nil
	}
	a.Key = prefix + a.Key
	return h.writeAttrFlat(a)
}

func (h *prettyHandler) writeAttrFlat(a slog.Attr) error {
	_, err := fmt.Fprintf(h.w, " %s=%s", a.Key, a.Value.String())
	return err
}

func (h *prettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)
	return &prettyHandler{
		opts:   h.opts,
		w:      h.w,
		mu:     h.mu,
		attrs:  newAttrs,
		groups: h.groups,
	}
}

func (h *prettyHandler) WithGroup(name string) slog.Handler {
	newGroups := make([]string, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = name
	return &prettyHandler{
		opts:   h.opts,
		w:      h.w,
		mu:     h.mu,
		attrs:  h.attrs,
		groups: newGroups,
	}
}

func levelToColor(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return colorRed
	case level >= slog.LevelWarn:
		return colorYellow
	case level >= slog.LevelInfo:
		return colorGreen
	default:
		return colorCyan
	}
}
