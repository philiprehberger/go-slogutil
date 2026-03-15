package slogutil

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestPrettyHandlerOutput(t *testing.T) {
	var buf bytes.Buffer
	h := PrettyHandler(&buf, &PrettyHandlerOptions{
		Level:      slog.LevelDebug,
		TimeFormat: time.DateTime,
	})

	r := slog.NewRecord(time.Date(2026, 3, 15, 10, 30, 0, 0, time.UTC), slog.LevelInfo, "hello world", 0)
	r.AddAttrs(slog.String("key", "value"))

	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatalf("Handle error: %v", err)
	}

	out := buf.String()

	if !strings.Contains(out, "2026-03-15") {
		t.Errorf("expected timestamp in output, got: %s", out)
	}
	if !strings.Contains(out, "INFO") {
		t.Errorf("expected level INFO in output, got: %s", out)
	}
	if !strings.Contains(out, "hello world") {
		t.Errorf("expected message in output, got: %s", out)
	}
	if !strings.Contains(out, "key=value") {
		t.Errorf("expected key=value in output, got: %s", out)
	}
}

func TestPrettyHandlerColorCodes(t *testing.T) {
	tests := []struct {
		level slog.Level
		color string
	}{
		{slog.LevelDebug, colorCyan},
		{slog.LevelInfo, colorGreen},
		{slog.LevelWarn, colorYellow},
		{slog.LevelError, colorRed},
	}

	for _, tt := range tests {
		var buf bytes.Buffer
		h := PrettyHandler(&buf, &PrettyHandlerOptions{Level: slog.LevelDebug})
		r := slog.NewRecord(time.Now(), tt.level, "test", 0)
		if err := h.Handle(context.Background(), r); err != nil {
			t.Fatalf("Handle error: %v", err)
		}
		if !strings.Contains(buf.String(), tt.color) {
			t.Errorf("expected color code %q for level %s", tt.color, tt.level)
		}
	}
}

func TestPrettyHandlerEnabled(t *testing.T) {
	h := PrettyHandler(nil, &PrettyHandlerOptions{Level: slog.LevelWarn})

	if h.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("expected Info to be disabled when min level is Warn")
	}
	if !h.Enabled(context.Background(), slog.LevelWarn) {
		t.Error("expected Warn to be enabled")
	}
	if !h.Enabled(context.Background(), slog.LevelError) {
		t.Error("expected Error to be enabled")
	}
}

func TestPrettyHandlerWithAttrs(t *testing.T) {
	var buf bytes.Buffer
	h := PrettyHandler(&buf, &PrettyHandlerOptions{Level: slog.LevelDebug})
	h2 := h.WithAttrs([]slog.Attr{slog.String("service", "api")})

	r := slog.NewRecord(time.Now(), slog.LevelInfo, "request", 0)
	if err := h2.Handle(context.Background(), r); err != nil {
		t.Fatalf("Handle error: %v", err)
	}

	if !strings.Contains(buf.String(), "service=api") {
		t.Errorf("expected pre-set attr in output, got: %s", buf.String())
	}
}

func TestPrettyHandlerWithGroup(t *testing.T) {
	var buf bytes.Buffer
	h := PrettyHandler(&buf, &PrettyHandlerOptions{Level: slog.LevelDebug})
	h2 := h.WithGroup("req")

	r := slog.NewRecord(time.Now(), slog.LevelInfo, "request", 0)
	r.AddAttrs(slog.String("method", "GET"))
	if err := h2.Handle(context.Background(), r); err != nil {
		t.Fatalf("Handle error: %v", err)
	}

	if !strings.Contains(buf.String(), "req.method=GET") {
		t.Errorf("expected grouped attr in output, got: %s", buf.String())
	}
}
