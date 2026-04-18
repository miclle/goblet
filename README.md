# Go Single Page Binary App

A project template for building Go + React single-page applications that compile into a single binary.

The backend embeds the frontend build output via `//go:embed`, so production deployment is typically a single executable plus a database.

## Tech Stack

- Backend: Go 1.25, `fox-gonic/fox`, GORM
- Database: PostgreSQL (default) / MySQL
- Frontend: React 18, TypeScript 5, Vite 6, Tailwind CSS 4, shadcn/ui
  - React Router v7
  - React Query v5

## Requirements

- Go 1.25+
- Node.js 22.14+
- PostgreSQL or MySQL
- [Task](https://taskfile.dev/)
- `reflex` for `task dev` hot reload
- `golangci-lint` for `task check`

## Quick Start

```bash
git clone https://github.com/miclle/go-single-page-binary-app.git
cd go-single-page-binary-app
task install
```

Create local config:

```bash
cp cmd/app/config.example.yaml cmd/app/config.local.yaml
# Edit config.local.yaml with your database settings
```

Start development:

```bash
task dev
```

## Common Commands

```bash
task install        # Install Go and frontend dependencies
task dev            # Start Vite + Go hot reload
task build          # Build the server binary with embedded frontend
task build-all      # Cross-build server binaries
task run            # Run the server with local config
task lint           # Auto-fix Go style and run frontend lint
task check          # CI-aligned checks without rewriting files
task test           # Go tests with race detection and coverage
task clean          # Remove build artifacts
task update-tools   # Install/update dev tools
```

## Architecture

```
cmd/app/              → Application entry point
internal/config/      → YAML config loading
internal/entity/      → Data models (GORM)
internal/handler/     → HTTP handlers + routes + middleware
internal/service/     → Business logic + DB operations
internal/errors/      → Centralized error types
internal/website/     → Embedded SPA (React + Vite)
pkg/gormlog/          → GORM logger adapter
```

### Single Binary Embedding

The key pattern: two Go files with build tags control how frontend assets are served:

- **Development** (`-tags development`): Reverse-proxies requests to Vite dev server at `localhost:5173`
- **Production** (default): Serves assets from `//go:embed build/*`, with SPA fallback for non-API routes

## Configuration

The YAML config keeps only bootstrap settings:

- `addr`: Listen address (e.g., `0.0.0.0:9000`)
- `driver`: Database driver (`postgres` or `mysql`)
- `dsn`: Database connection string

## Build & Deployment

```bash
task build          # Build production binary (includes frontend)
./bin/app -c config.yaml
```

## License

[Apache-2.0](https://www.apache.org/licenses/LICENSE-2.0)
