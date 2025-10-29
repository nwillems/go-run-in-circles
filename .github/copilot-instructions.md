# Copilot Instructions for go-run-in-circles

## Project Overview

This is a **lap tracking system** for running events. Current state: basic Go web server with embedded static files. The README.md contains the complete architecture specification for the target system.

## Architecture Pattern

- **Embedded Web App**: Go uses `//go:embed web/*` to bundle HTML/JS/CSS into the binary
- **Single Binary Deployment**: No external dependencies, just the compiled Go binary + SQLite file
- **File Structure**: All Go code in `./`, web assets in `./web/`

## Key Implementation Details

### Current State vs Target

**Current**: Basic file server serving static content from `./web/`
**Target**: Full lap tracking system with SQLite backend (see README.md for complete spec)

### Golang coding style

- Idiomatic Go: clear, concise, error handling
- Use standard library where possible
- Modular functions for handlers, DB access, utils

If a golang file exceeds ~100 lines, consider breaking it into smaller files.

### Go Web Server Pattern

```go
//go:embed web/*
var content embed.FS

func main() {
    http.Handle("/", http.FileServer(http.FS(content)))
    // Add API routes here: /runners, /laps, /dashboard
}
```

### Required API Endpoints (from README.md)

- `POST /runners` - Register runner (bib, name, photo as BLOB)
- `POST /laps` - Record lap (bib number, timestamp)
- `GET /dashboard` - Return all runners with calculated metrics

### Database Schema (SQLite)

- `runners`: bib_number (PK), name, photo (BLOB)
- `laps`: id (PK), bib_number (FK), lap_number, timestamp
- `events`: id (PK), name, date, lap_distance

### Frontend Architecture

- **Auto-refresh**: Poll `/dashboard` every 5 seconds
- **Photo handling**: Upload as base64, store as BLOB in SQLite
- **Two views**: Official dashboard + public display screen

## Development Workflows

### Running the Server

```bash
go build .
./go-run-in-circles event.db
# Serves on http://localhost:8080
```

## Critical Patterns

1. **Portability Focus**: Single binary + database file, no external services
2. **Photo Storage**: Base64 frontend → BLOB database (not file system)
3. **Real-time Updates**: Frontend polling, not WebSockets
4. **Database Selection**: User selects SQLite file at startup (no reset functionality)

## Next Implementation Steps

1. Add SQLite integration with the three-table schema
2. Implement the three API endpoints with proper JSON responses
3. Build frontend with runner registration and lap recording forms
4. Add the auto-refresh dashboard with calculated metrics (avg lap time, pace, etc.)
