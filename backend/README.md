# Calculator Backend (Go)

A small, dependency-free REST microservice exposing arithmetic operations.
Built with the Go standard library (`net/http`) — no web framework.

- **Go:** 1.26.6
- **Version:** 0.0.1 (SemVer)

## Layout

```
cmd/server        Composition root: config → wiring → HTTP server
internal/
  calculator      Pure arithmetic domain (registry, errors) — no HTTP
  api             Transport: DTOs, handlers, router, middleware
  config          Environment configuration with defaults
  version         SemVer constant
```

The domain (`internal/calculator`) knows nothing about HTTP; the transport
(`internal/api`) depends on a narrow `Calculator` interface, so each layer is
tested in isolation.

## Run

```bash
go run ./cmd/server            # listens on :8080
PORT=9090 go run ./cmd/server  # custom port
```

Configuration (all optional):

| Variable      | Default | Description                          |
| ------------- | ------- | ------------------------------------ |
| `PORT`        | `8080`  | TCP port to listen on                |
| `CORS_ORIGIN` | `*`     | Comma-separated allowed CORS origins |

## Test & coverage

```bash
go test ./...                                   # run tests
go test ./... -coverprofile=coverage.out        # with coverage
go tool cover -func=coverage.out                # summary in terminal
go tool cover -html=coverage.out -o coverage.html  # HTML report
```

Domain and config packages are covered 100%; the transport layer ~97%. Only the
`main` composition root is intentionally left uncovered.

## API

See the [root README](../README.md#api-reference) for the full API reference,
request/response schemas, and `curl` examples.
