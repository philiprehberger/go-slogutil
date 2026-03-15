package slogutil

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestSamplingRate(t *testing.T) {
	var buf bytes.Buffer
	inner := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	h := Sampling(inner, 10)

	for i := 0; i < 100; i++ {
		r := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
		if err := h.Handle(context.Background(), r); err != nil {
			t.Fatalf("Handle error: %v", err)
		}
	}

	// With rate=10, exactly 10 out of 100 should be logged (counter 10,20,...,100).
	lines := strings.Count(buf.String(), "\n")
	if lines != 10 {
		t.Errorf("expected 10 sampled records, got %d", lines)
	}
}

func TestSamplingAlwaysLogsErrors(t *testing.T) {
	var buf bytes.Buffer
	inner := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	h := Sampling(inner, 1000) // Very high rate, almost nothing sampled.

	for i := 0; i < 50; i++ {
		r := slog.NewRecord(time.Now(), slog.LevelError, "err", 0)
		if err := h.Handle(context.Background(), r); err != nil {
			t.Fatalf("Handle error: %v", err)
		}
	}

	lines := strings.Count(buf.String(), "err")
	if lines != 50 {
		t.Errorf("expected all 50 error records to be logged, got %d", lines)
	}
}

func TestSamplingRateOne(t *testing.T) {
	var buf bytes.Buffer
	inner := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	h := Sampling(inner, 1) // rate <= 1 returns inner directly.

	r := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatalf("Handle error: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("expected record to be logged with rate=1")
	}
}
