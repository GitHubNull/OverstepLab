> ⚠️ **SECURITY WARNING**: OverstepLab is a **deliberately vulnerable application** intended solely for educational purposes. Do NOT deploy to production or expose to the public internet. See [SECURITY.md](SECURITY.md).

# OverstepLab - Privilege Escalation Testing Range

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org/)
[![Vue](https://img.shields.io/badge/Vue-3.4-4FC08D.svg)](https://vuejs.org/)

**[中文文档](README.md)**

## Overview

OverstepLab is a **privilege escalation testing range** simulating a multi-user VPS cloud management platform (CloudNest). It contains intentionally implanted access control vulnerabilities for security professionals to practice vulnerability discovery and secure coding.

## Features

- 🎯 **13 Vulnerability Scenarios**: Horizontal IDOR, Vertical Escalation, Context-based Escalation
- 🏢 **Realistic Business Scenario**: Multi-user VPS platform with company/individual accounts
- 🔄 **Security/Vulnerable Mode Toggle**: Switch between modes for comparison
- 💡 **Progressive Hint System**: 3 levels of hints per vulnerability
- 📋 **Complete WriteUps**: Detailed vulnerability analysis and remediation
- 🔧 **One-Click Database Reset**: Restore initial state anytime
- 📦 **Single Binary Deployment**: Go compiles with embedded frontend assets
- 🚀 **Cross-Platform One-Click Scripts**: Windows (.bat/.ps1), Linux/macOS (.sh) double-click to run
- 🐳 **Docker Support**: One-command container deployment

## Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.21+ / Gin / GORM / SQLite (pure Go) / JWT |
| Frontend | Vue 3 / PrimeVue 4 / Pinia / Vite |
| Database | SQLite3 (embedded) |

## Quick Start

### Option 1: One-Click Startup Scripts (Recommended)

The easiest way to start. Automatically builds frontend, compiles backend, and starts the server.

**Windows:**
```powershell
# Double-click to run
start.bat
# or
start.ps1
```

**Linux / macOS:**
```bash
# Grant execute permission then run
chmod +x start.sh
./start.sh
```

Access `http://localhost:5000` after startup.

### Option 2: Development Mode

Requires two terminals for frontend and backend.

```bash
# Terminal 1: Start backend
cd src/backend && go mod tidy && go run ./cmd/server/main.go

# Terminal 2: Start frontend
cd src/frontend && pnpm install && pnpm dev
```

### Option 3: Production (Single Binary)

```bash
make build
./bin/oversteplab
# Access http://localhost:5000
```

### Option 4: Docker

```bash
docker-compose up --build
# Access http://localhost:5000
```

## Pre-seeded Test Accounts

| Username | Password | Type | Role | Company |
|----------|----------|------|------|---------|
| admin | admin123 | Platform Admin | platform_admin | - |
| acme_admin | pass123 | Company | admin | Acme Corp |
| acme_ops | pass123 | Company | operator | Acme Corp |
| acme_finance | pass123 | Company | finance | Acme Corp |
| acme_viewer | pass123 | Company | viewer | Acme Corp |
| globex_admin | pass123 | Company | admin | Globex Inc |
| globex_ops | pass123 | Company | operator | Globex Inc |
| globex_finance | pass123 | Company | finance | Globex Inc |
| globex_viewer | pass123 | Company | viewer | Globex Inc |
| alice | pass123 | Individual | - | - |
| bob | pass123 | Individual | - | - |

## Vulnerability List

| ID | Category | Scenario | Difficulty |
|----|----------|----------|------------|
| H-01 | Horizontal IDOR | View other user's VPS by modifying ID | ★☆☆ |
| H-02 | Horizontal IDOR | View other user's profile by modifying ID | ★☆☆ |
| H-03 | Horizontal IDOR | View other user's order by modifying ID | ★☆☆ |
| H-04 | Horizontal IDOR | View/reply other user's ticket by modifying ID | ★☆☆ |
| H-05 | Horizontal IDOR | Delete other user's API Key by modifying ID | ★★☆ |
| V-01 | Vertical | Viewer role controls VPS power | ★☆☆ |
| V-02 | Vertical | Operator calls user management API to add member | ★★☆ |
| V-03 | Vertical | Finance calls VPS reinstall API | ★★☆ |
| V-04 | Vertical | Individual user calls platform admin API | ★★☆ |
| V-05 | Vertical | Operator upgrades self role to admin | ★★★ |
| C-01 | Context | Cross-company VPS operation | ★★☆ |
| C-02 | Context | Individual user creates company member | ★★★ |
| C-03 | Context | Revoked API Key still accessible | ★★★ |

See full details in `doc/tutorials/06-writeups.md`.

## License

MIT License. See [LICENSE](LICENSE).

## Security

This project contains **intentionally implanted vulnerabilities** for educational purposes only. See [SECURITY.md](SECURITY.md).
