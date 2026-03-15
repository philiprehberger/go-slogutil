package slogutil

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Duration returns a slog.Attr with key "duration" and a human-readable
// string representation of d (e.g. "1.5s", "200ms").
func Duration(d time.Duration) slog.Attr {
	return slog.String("duration", d.String())
}

// Error returns a slog.Attr with key "error" and the error message as value.
// If err is nil, the attr value is "<nil>".
func Error(err error) slog.Attr {
	if err == nil {
		return slog.String("error", "<nil>")
	}
	return slog.String("error", err.Error())
}

// HTTPRequest returns a slog.Attr group named "http_request" containing the
// method, path, and remote address from the given request.
func HTTPRequest(r *http.Request) slog.Attr {
	return slog.Group("http_request",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("remote_addr", r.RemoteAddr),
	)
}

// HTTPStatus returns a slog.Attr with key "status" and the HTTP status code
// as an integer value.
func HTTPStatus(code int) slog.Attr {
	return slog.Int("status", code)
}

// TraceID returns a slog.Attr with key "trace_id" and the given trace
// identifier as a string value.
func TraceID(id string) slog.Attr {
	return slog.String("trace_id", id)
}

// Stringer returns a slog.Attr with the given key and the result of calling
// String() on v. This is a convenience for types that implement fmt.Stringer.
func Stringer(key string, v fmt.Stringer) slog.Attr {
	return slog.String(key, v.String())
}
