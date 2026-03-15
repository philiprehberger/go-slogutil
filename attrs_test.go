package slogutil

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	attr := Duration(500 * time.Millisecond)
	if attr.Key != "duration" {
		t.Errorf("expected key 'duration', got %q", attr.Key)
	}
	if attr.Value.String() != "500ms" {
		t.Errorf("expected value '500ms', got %q", attr.Value.String())
	}
}

func TestError(t *testing.T) {
	attr := Error(errors.New("something failed"))
	if attr.Key != "error" {
		t.Errorf("expected key 'error', got %q", attr.Key)
	}
	if attr.Value.String() != "something failed" {
		t.Errorf("expected value 'something failed', got %q", attr.Value.String())
	}
}

func TestErrorNil(t *testing.T) {
	attr := Error(nil)
	if attr.Key != "error" {
		t.Errorf("expected key 'error', got %q", attr.Key)
	}
	if attr.Value.String() != "<nil>" {
		t.Errorf("expected value '<nil>', got %q", attr.Value.String())
	}
}

func TestHTTPRequest(t *testing.T) {
	req := &http.Request{
		Method:     "GET",
		URL:        &url.URL{Path: "/api/users"},
		RemoteAddr: "127.0.0.1:8080",
	}
	attr := HTTPRequest(req)
	if attr.Key != "http_request" {
		t.Errorf("expected key 'http_request', got %q", attr.Key)
	}
	if attr.Value.Kind() != slog.KindGroup {
		t.Fatalf("expected group kind, got %v", attr.Value.Kind())
	}

	group := attr.Value.Group()
	found := map[string]string{}
	for _, a := range group {
		found[a.Key] = a.Value.String()
	}

	if found["method"] != "GET" {
		t.Errorf("expected method=GET, got %q", found["method"])
	}
	if found["path"] != "/api/users" {
		t.Errorf("expected path=/api/users, got %q", found["path"])
	}
	if found["remote_addr"] != "127.0.0.1:8080" {
		t.Errorf("expected remote_addr=127.0.0.1:8080, got %q", found["remote_addr"])
	}
}

func TestHTTPStatus(t *testing.T) {
	attr := HTTPStatus(200)
	if attr.Key != "status" {
		t.Errorf("expected key 'status', got %q", attr.Key)
	}
	if attr.Value.Int64() != 200 {
		t.Errorf("expected value 200, got %d", attr.Value.Int64())
	}
}

func TestTraceID(t *testing.T) {
	attr := TraceID("abc-123-def")
	if attr.Key != "trace_id" {
		t.Errorf("expected key 'trace_id', got %q", attr.Key)
	}
	if attr.Value.String() != "abc-123-def" {
		t.Errorf("expected value 'abc-123-def', got %q", attr.Value.String())
	}
}
