# AGENTS.md - AI Agent Guide for OverstepLab

## Project Identity

OverstepLab is a deliberately vulnerable privilege escalation testing range simulating a VPS management platform. Every vulnerability is intentional and must remain discoverable.

## Architecture

### Backend (Go)
- **Framework**: Gin (HTTP router)
- **ORM**: GORM
- **Database**: SQLite3 via `github.com/glebarez/sqlite` (pure Go)
- **Auth**: JWT via `github.com/golang-jwt/jwt/v5`
- **Structure**:
  - `internal/handler/` - HTTP handlers (controllers)
  - `internal/service/` - Business logic with vulnerability branching
  - `internal/repository/` - Data access layer
  - `internal/model/` - GORM structs
  - `internal/middleware/` - JWT auth, RBAC, CORS, audit logging
  - `internal/vuln/` - Security mode toggle and challenge metadata
  - `internal/web/` - Embedded frontend assets

### Frontend (Vue 3)
- **Framework**: Vue 3 Composition API
- **UI**: PrimeVue 4
- **State**: Pinia
- **Build**: Vite
- **Structure**:
  - `src/api/` - Axios wrappers
  - `src/stores/` - Pinia stores
  - `src/views/` - Page components
  - `src/router/` - Vue Router with auth guards

### Key Files
- `src/backend/cmd/server/main.go` - Entry point, DB init, server start
- `src/backend/internal/vuln/mode.go` - Security mode toggle
- `src/backend/internal/vuln/challenges.go` - Challenge metadata (13 vulnerabilities)
- `src/backend/internal/web/embed.go` - `//go:embed` directive for frontend

## Coding Conventions

### Go
- Layered architecture: handler → service → repository → model
- Return DTOs, not models
- Vulnerabilities implemented as mode-conditional branching in service layer

### Vue
- `<script setup>` Composition API
- Pinia for global state
- `v-if` for UI-level permission gating (intentionally bypassable)

## Intentional Vulnerability Protocol

- All vulnerabilities catalogued in `doc/PRD.md` §4.2 and `doc/tutorials/06-writeups.md`
- **NEVER "fix" a vulnerability** unless explicitly tasked to implement Safe Mode
- New vulnerability additions must update: PRD, tutorials, challenge metadata
- Safe mode controlled by `OVERSTEPLAB_SAFE_MODE` env var and runtime toggle

## Safe Mode

- `vuln.IsSecureMode()` returns current state
- When true, all vulnerabilities are patched via service-layer checks
- When false (default), vulnerabilities are exploitable

## Development Workflow

```bash
# Quick start (one command)
./start.sh         # Linux/macOS
start.bat          # Windows CMD
start.ps1          # Windows PowerShell

# Dev mode (two terminals)
make dev-frontend  # Vite on :5173
make dev-backend   # Go on :5000

# Production build
make build         # Frontend → embed dir → Go binary
```
