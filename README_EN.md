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

Access `http://localhost:8080` after startup.

### Option 2: Development Mode

Requires two terminals for frontend and backend.

```bash
# Terminal 1: Start backend
cd src/backend && go run ./cmd/server/main.go

# Terminal 2: Start frontend
cd src/frontend && pnpm install && pnpm dev
```

### Option 3: Production (Single Binary)

```bash
make build
./bin/oversteplab
# Access http://localhost:8080
```

### Option 4: Docker

```bash
docker-compose up --build
# Access http://localhost:8080
```

## Pre-seeded Test Accounts

| Username | Password | Type | Role | Company |
|----------|----------|------|------|---------|
| admin | admin123 | Platform Admin | platform_admin | - |
| acme_admin | pass123 | Company | admin | Acme Corp |
| acme_ops | pass123 | Company | operator | Acme Corp |
| alice | pass123 | Individual | - | - |

## Vulnerability List

| ID | Category | Scenario | Difficulty |
|----|----------|----------|------------|
| H-01 | Horizontal IDOR | View other user's VPS by modifying ID | ★☆☆ |
| V-01 | Vertical | Viewer role controls VPS power | ★☆☆ |
| C-01 | Context | Cross-company VPS operation | ★★☆ |

See full list in [README.md](README.md) or `doc/tutorials/06-writeups.md`.

## License

MIT License. See [LICENSE](LICENSE).

## Security

This project contains **intentionally implanted vulnerabilities** for educational purposes only. See [SECURITY.md](SECURITY.md).
