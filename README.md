# pitest3

A Todo REST API server and CLI built in Go.

## Quick Start

```bash
# Start the server (default port 8080)
go run ./cmd/server

# Start on a custom port
go run ./cmd/server --port 9090

# Use a custom data file
go run ./cmd/server --data-file /path/to/data.json
```

## Project Layout

```
cmd/server       — HTTP server entry point
cmd/todos        — CLI entry point
internal/model   — domain models
internal/storage — persistence layer
internal/handler — HTTP handlers
```

## Development

```bash
go build ./...   # build all packages
go test ./...    # run tests
go vet ./...     # static analysis
```

## Glossary

See [CONTEXT.md](CONTEXT.md) for domain terminology.
