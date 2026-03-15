package slogutil

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestLevelRouteExactMatch(t *testing.T) {
	var infoBuf, errorBuf bytes.Buffer
	infoH := slog.NewTextHandler(&infoBuf, &slog.HandlerOptions{Level: slog.LevelDebug})
	errorH := slog.NewTextHandler(&errorBuf, &slog.HandlerOptions{Level: slog.LevelDebug})

	router := LevelRoute(map[slog.Level]slog.Handler{
		slog.LevelInfo:  infoH,
		slog.LevelError: errorH,
	})

	infoRec := slog.NewRecord(time.Now(), slog.LevelInfo, "info msg", 0)
	if err := router.Handle(context.Background(), infoRec); err != nil {
		t.Fatalf("Handle error: %v", err)
	}

	errRec := slog.NewRecord(time.Now(), slog.LevelError, "error msg", 0)
	if err := router.Handle(context.Background(), errRec); err != nil {
		t.Fatalf("Handle error: %v", err)
	}

	if infoBuf.Len() == 0 {
		t.Error("expected info handler to receive info record")
	}
	if errorBuf.Len() == 0 {
		t.Error("expected error handler to receive error record")
	}
}

func TestLevelRouteFallback(t *testing.T) {
	var infoBuf bytes.Buffer
	infoH := slog.NewTextHandler(&infoBuf, &slog.HandlerOptions{Level: slog.LevelDebug})

	router := LevelRoute(map[slog.Level]slog.Handler{
		slog.LevelInfo: infoH,
	})

	// Warn should fall back to Info handler since Warn > Info and no Warn handler.
	warnRec := slog.NewRecord(time.Now(), slog.LevelWarn, "warn msg", 0)
	if err := router.Handle(context.Background(), warnRec); err != nil {
		t.Fatalf("Handle error: %v", err)
	}

	if infoBuf.Len() == 0 {
		t.Error("expected info handler to receive warn record via fallback")
	}
}

func TestLevelRouteNoMatch(t *testing.T) {
	var errorBuf bytes.Buffer
	errorH := slog.NewTextHandler(&errorBuf, &slog.HandlerOptions{Level: slog.LevelDebug})

	router := LevelRoute(map[slog.Level]slog.Handler{
		slog.LevelError: errorH,
	})

	// Debug has no handler and no lower-level fallback.
	debugRec := slog.NewRecord(time.Now(), slog.LevelDebug, "debug msg", 0)
	if err := router.Handle(context.Background(), debugRec); err != nil {
		t.Fatalf("Handle error: %v", err)
	}

	if errorBuf.Len() != 0 {
		t.Error("expected no handler to receive debug record")
	}
}

func TestLevelRouteEnabled(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})

	router := LevelRoute(map[slog.Level]slog.Handler{
		slog.LevelInfo: h,
	})

	if !router.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("expected enabled for info level")
	}
}
