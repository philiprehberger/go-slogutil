package slogutil

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestMultiFanOut(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	h1 := slog.NewTextHandler(&buf1, &slog.HandlerOptions{Level: slog.LevelDebug})
	h2 := slog.NewTextHandler(&buf2, &slog.HandlerOptions{Level: slog.LevelDebug})
	m := Multi(h1, h2)

	r := slog.NewRecord(time.Now(), slog.LevelInfo, "hello", 0)
	if err := m.Handle(context.Background(), r); err != nil {
		t.Fatalf("Handle error: %v", err)
	}

	if buf1.Len() == 0 {
		t.Error("expected handler 1 to receive the record")
	}
	if buf2.Len() == 0 {
		t.Error("expected handler 2 to receive the record")
	}
}

func TestMultiEnabled(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	h1 := slog.NewTextHandler(&buf1, &slog.HandlerOptions{Level: slog.LevelError})
	h2 := slog.NewTextHandler(&buf2, &slog.HandlerOptions{Level: slog.LevelDebug})
	m := Multi(h1, h2)

	// Should be enabled because h2 accepts Debug.
	if !m.Enabled(context.Background(), slog.LevelDebug) {
		t.Error("expected Enabled=true when at least one handler accepts the level")
	}
}

func TestMultiEnabledNone(t *testing.T) {
	var buf1 bytes.Buffer
	h1 := slog.NewTextHandler(&buf1, &slog.HandlerOptions{Level: slog.LevelError})
	m := Multi(h1)

	if m.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("expected Enabled=false when no handler accepts the level")
	}
}

func TestMultiWithAttrs(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	h1 := slog.NewTextHandler(&buf1, &slog.HandlerOptions{Level: slog.LevelDebug})
	h2 := slog.NewTextHandler(&buf2, &slog.HandlerOptions{Level: slog.LevelDebug})
	m := Multi(h1, h2).WithAttrs([]slog.Attr{slog.String("k", "v")})

	r := slog.NewRecord(time.Now(), slog.LevelInfo, "test", 0)
	if err := m.Handle(context.Background(), r); err != nil {
		t.Fatalf("Handle error: %v", err)
	}

	if buf1.Len() == 0 || buf2.Len() == 0 {
		t.Error("expected both handlers to receive the record with attrs")
	}
}
