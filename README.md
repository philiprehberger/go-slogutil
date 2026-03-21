# go-slogutil

[![CI](https://github.com/philiprehberger/go-slogutil/actions/workflows/ci.yml/badge.svg)](https://github.com/philiprehberger/go-slogutil/actions/workflows/ci.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/philiprehberger/go-slogutil.svg)](https://pkg.go.dev/github.com/philiprehberger/go-slogutil) [![License](https://img.shields.io/github/license/philiprehberger/go-slogutil)](LICENSE)

Handlers, formatters, and helpers for Go's log/slog. The missing batteries

## Installation

```bash
go get github.com/philiprehberger/go-slogutil
```

## Usage

### Pretty Console Output (Development)

```go
import "github.com/philiprehberger/go-slogutil"

logger := slog.New(slogutil.PrettyHandler(os.Stderr, &slogutil.PrettyHandlerOptions{
    Level:      slog.LevelDebug,
    TimeFormat: time.DateTime,
}))
logger.Info("server started", "port", 8080)
// Output: 2026-03-15 10:30:00 INFO server started port=8080
```

### Multi Handler (Fan-out)

```go
console := slogutil.PrettyHandler(os.Stderr, nil)
file := slog.NewJSONHandler(logFile, nil)

logger := slog.New(slogutil.Multi(console, file))
logger.Info("logged to both console and file")
```

### Level Routing

```go
stdout := slog.NewTextHandler(os.Stdout, nil)
stderr := slog.NewTextHandler(os.Stderr, nil)

logger := slog.New(slogutil.LevelRoute(map[slog.Level]slog.Handler{
    slog.LevelInfo:  stdout,
    slog.LevelError: stderr,
}))
logger.Info("goes to stdout")
logger.Error("goes to stderr")
```

### Sampling

```go
// Log only 1 in every 100 info messages (errors are always logged).
h := slogutil.Sampling(slog.NewJSONHandler(os.Stdout, nil), 100)
logger := slog.New(h)
```

### Attribute Helpers

```go
logger.Info("request handled",
    slogutil.Duration(elapsed),
    slogutil.HTTPRequest(r),
    slogutil.HTTPStatus(200),
    slogutil.TraceID("abc-123"),
)

if err != nil {
    logger.Error("failed", slogutil.Error(err))
}
```

### Domain Attributes

```go
import slogutil "github.com/philiprehberger/go-slogutil"

logger.Info("query complete",
    slogutil.Database("users", "localhost", 5432),
    slogutil.UserID("user-123"),
    slogutil.RequestID("req-abc"),
    slogutil.Latency(42 * time.Millisecond),
)
```

## API

| Function / Type | Description |
|-----------------|-------------|
| `PrettyHandler(w, opts)` | Colorized console handler for development |
| `PrettyHandlerOptions` | Options: Level, TimeFormat |
| `Multi(handlers...)` | Fan-out to multiple handlers |
| `LevelRoute(routes)` | Route records by level to different handlers |
| `Sampling(handler, rate)` | Log 1 in every N records (errors always pass) |
| `Duration(d)` | Formatted duration attr |
| `Error(err)` | Error attr with "error" key |
| `HTTPRequest(r)` | Group attr with method, path, remote_addr |
| `HTTPStatus(code)` | Status code attr |
| `TraceID(id)` | Trace ID attr |
| `Stringer(key, v)` | Attr from any fmt.Stringer |
| `Database(name, host, port)` | Group attr with database connection details |
| `UserID(id)` | User identifier attr |
| `RequestID(id)` | Request identifier attr |
| `Latency(d)` | Latency in milliseconds attr |

## Development

```bash
go test ./...
go vet ./...
```

## License

MIT
