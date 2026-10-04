# go-monolith

A Go backend API built from scratch using only the standard library (`net/http`), focused on learning idiomatic Go and production-oriented backend engineering practices.

## Overview

This project implements a listings API as a vehicle for practicing clean backend architecture in Go. The emphasis is not on the domain logic itself, but on the engineering decisions behind it: how dependencies are wired, how context flows through the stack, how errors are surfaced to clients, and how the database layer is managed responsibly.

Deployed to [Render](https://go-monolith-w6pr.onrender.com) from day 1 -- shipping incrementally instead of building everything locally and deploying at the end. No frameworks. No routers. Just Go 1.22+ `net/http` with method-based routing, a live Neon PostgreSQL instance, and deliberate software design.

## Engineering Practices

**Dependency Injection** -- Handlers get their dependencies (`*sql.DB`, `*slog.Logger`) passed in through constructors, not globals. No DI framework needed.

**Context Propagation** -- Every handler pulls `r.Context()` and passes it into `QueryContext` / `ExecContext` / `QueryRowContext`. If a client disconnects or a timeout hits, the database query gets cancelled instead of running as a zombie.

**Structured Logging** -- Uses `log/slog` with JSON output. Logs include structured fields like `request_id`, `listing_id`, and `err` instead of plain string messages.

**Request ID Middleware** -- Generates or forwards `X-Request-ID` headers, stores them in context using an unexported key type (avoids collisions with other packages), and threads them into log entries for tracing.

**Standardized Error Responses** -- All API errors return a consistent JSON shape: `{ "error": { "code", "message", "field?" } }`. Error codes are typed constants (`type Code string`) so they can't be mixed up with regular strings.

**DTOs and Validation** -- Request and response structs are separate from database models. Validation errors include the field name so the client knows exactly what went wrong.

**Connection Pool Tuning** -- The database pool has explicit limits (`MaxOpenConns`, `MaxIdleConns`, `ConnMaxLifetime`) and pings the database at startup with a timeout to fail fast.

**Server Timeouts** -- `http.Server` is configured with read, write, and idle timeouts to protect against slow or hanging clients.

**Database Migrations** -- Schema changes use `golang-migrate` with versioned up/down SQL files and a separate `cmd/migrate` entrypoint.

**Environment Config** -- Config loads from environment variables (`.env` via `godotenv`), validated at startup. Missing values panic immediately instead of silently defaulting.

## Architecture

```
cmd/
  api/          -- Application entrypoint: wires dependencies, configures server
  migrate/      -- Migration CLI: applies or rolls back schema changes

internal/
  config/       -- Environment-based configuration loading and validation
  db/           -- Connection pool setup with health check
  handlers/     -- HTTP handlers, DTOs, request validation
  httpx/        -- Standardized JSON error envelope utilities
  middleware/   -- Request ID generation and context propagation

migrations/     -- Versioned SQL migration files (up/down)
```

## Tech Stack

- **Go 1.22+** -- standard library `net/http` with method-based routing (no third-party router)
- **PostgreSQL** -- via `pgx` driver through `database/sql`
- **golang-migrate** -- versioned schema migrations
- **log/slog** -- structured JSON logging
- **godotenv** -- environment variable management

## Running Locally

```bash
# Set up environment
cp .env.example .env   # Configure PORT, ENV, DATABASE_URL

# Run migrations
make migrate-up

# Build and run
make run
```

The API starts on the configured port with endpoints:

```
GET    /healthz          -- Health check
GET    /listings          -- List all listings
POST   /listings          -- Create a listing
DELETE /listings/{id}     -- Delete a listing
```

## What This Demonstrates

This project is a focused exercise in writing Go the way it's meant to be written: explicit error handling, context-aware database operations, constructor-based dependency injection, structured observability, and clear separation between HTTP transport, business logic, and data access. Every design decision here reflects a specific lesson in building backends that are maintainable, debuggable, and safe under real-world conditions.

