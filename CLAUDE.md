# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Go Core is a collection of common packages for building and deploying Go applications. The repository provides reusable utilities for database access, logging, web servers, cryptography, context management, and environment configuration.

## Development Commands

### Initial Setup
```bash
make install_deps  # Install all dependencies including golangci-lint, goimports, goverreport, and pre-commit
```

### Testing
```bash
# Run full test suite with coverage (90% threshold required)
make coverage_tests

# The coverage test process:
# 1. Starts test databases via docker-compose (MySQL, PostgreSQL, SQLite)
# 2. Runs all tests with coverage profiling
# 3. Shuts down test databases
# 4. Generates coverage report and fails if below 90%
# 5. Opens HTML coverage report in browser
```

### Linting
```bash
# Format and lint code
make lint

# This runs:
# - gofmt for formatting
# - goimports for import organization
# - golangci-lint with custom config (.golangci.yaml)
# - Only checks changes from HEAD~1 (new code)
```

### Running Individual Tests
```bash
# Run tests for a specific package
go test ./pkg/godb/... -v

# Run a specific test function
go test ./pkg/godb/... -run TestDatabaseConnection -v

# Run with coverage for a single package
go test ./pkg/golog/... -coverprofile=coverage.out
```

## Architecture

### Package Structure

The codebase follows a modular package design where each package in `pkg/` is independent and provides specific functionality:

**pkg/godb** - Database abstraction layer
- Wraps `sqlx` and provides typed interfaces for SQL operations
- Supports MySQL, PostgreSQL, and SQLite3
- Custom wrapper types (DB, Tx, Conn, Rows, Stmt, NamedStmt) that extend sqlx functionality
- Each interface method returns custom types instead of sqlx types
- Uses `customDB`, `customTx`, etc. structs that wrap underlying sqlx types
- Test error injection system via `pushTestError`/`popTestError` for testing error paths
- Connection timeout handling with context-based timeouts
- DSN building for different database types with default parameters

**pkg/golog** - Structured logging with Datadog integration
- Environment-aware logging (Test/Local uses colored console output, other envs use JSON with zap)
- Automatic Datadog trace/span ID injection when tracing context is available
- Uses `golog.WithTags()` to add structured key-value context to logs
- Four log levels: Info, Warn, Debug, Error
- Singleton pattern with `sync.Once` initialization

**pkg/goweb** - Fiber-based web server
- Built on `gofiber/fiber/v2` framework
- `DefaultConfig()` sets up standard middleware: CORS, rate limiting, panic recovery, favicon, profiling
- Automatic Swagger UI integration at configurable route
- Custom JSON encoder/decoder support (useful for sonic or other fast JSON libraries)
- Graceful shutdown handling via OS signals
- Routes defined via `WebRoute` struct with method, path, and handlers

**pkg/gocontext** - Type-safe context value storage
- Wrapper around `sync.Map` stored in context
- Generic `Add[T]` and `Get[T]` functions with type constraints
- Supports integer, float, string, and bool types
- Thread-safe concurrent access

**pkg/gocrypto** - Password hashing utilities
- Argon2id-based password hashing
- Configurable memory, iterations, parallelism, salt length, and key length
- Encoded hash format: `argon2id$version$memory$iterations$parallelism$salt$hash`
- Constant-time comparison to prevent timing attacks

**pkg/goenv** - Environment variable loader
- Generic `Load[T]()` function with type inference
- Supports int64, float64, string, bool types
- Returns default value if env var not set
- Panics on parse errors for fail-fast behavior

### Key Architectural Patterns

**Interface Wrapping Pattern (godb)**
The godb package wraps all sqlx types to provide a consistent interface while maintaining full sqlx compatibility. This allows:
- Mock-friendly testing (all types are interfaces)
- Custom error injection for testing
- Potential future enhancements without breaking API
- Type safety and IDE autocomplete for all database operations

**Environment-Based Initialization (golog)**
The logger automatically adjusts output format and verbosity based on environment:
- Test/Local: Human-readable colored console logs
- Development/Staging: JSON logs with debug level
- Production: JSON logs with warn level, no debug output

**Datadog Integration**
The logger automatically injects Datadog trace and span IDs from context when available via `gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer`. This enables correlation between logs and distributed traces.

## Testing Practices

- Test files use `_test.go` suffix
- Mock files use `_mock.go` or `_mock_test.go` suffix (mock_test.go files are only compiled during tests)
- Database tests require running databases (started via `make pre_coverage_tests`)
- 90% code coverage threshold enforced
- Use `testify/assert` for assertions
- Mock implementations follow same interface as real implementations

## Pull Request Requirements

According to `.github/pull_request_template.md`:

Mandatory checklist:
- Changes must be tested
- Functional and coverage tests added with real test values
- CHANGELOG (changelog.md) must be updated

Conditional:
- Add metrics for monitoring if the change affects runtime behavior
- Generate documentation if applicable

## Go Version and Dependencies

- Go 1.25.0
- Uses `bytedance/sonic` for fast JSON encoding (if configured in goweb)
- DataDog tracing: `gopkg.in/DataDog/dd-trace-go.v1`
- Fiber v2 for HTTP server
- sqlx with support for MySQL, PostgreSQL, SQLite3
- Zap for structured logging
- Argon2 for password hashing

## Linting Configuration

The project uses golangci-lint with custom configuration (`.golangci.yaml`):
- Enabled linters: asciicheck, bodyclose, dupl, errname, cyclop
- Maximum cyclomatic complexity: 15
- Tests are skipped for complexity checks
- Only checks new code from HEAD~1 to reduce noise
