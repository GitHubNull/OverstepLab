> ⚠️ **SECURITY WARNING**: OverstepLab is a **deliberately vulnerable application** intended solely for educational purposes. Do NOT deploy to production or expose to the public internet. See [SECURITY.md](SECURITY.md).

# OverstepLab - Privilege Escalation Testing Range

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8.svg)](https://golang.org/)
[![Vue](https://img.shields.io/badge/Vue-3.4-4FC08D.svg)](https://vuejs.org/)

**[中文文档](README.md)**

## Overview

OverstepLab is a **privilege escalation testing range** simulating a multi-user VPS cloud management platform (CloudNest). It contains 22 intentionally implanted access control vulnerabilities (13 traditional + 9 encoding/crypto challenges) for security professionals to practice vulnerability discovery and secure coding.

## Features

- 🎯 **22 Vulnerability Scenarios**:
  - 13 traditional privilege escalation vulnerabilities (Horizontal IDOR, Vertical Escalation, Context-based Escalation)
  - 9 encoding/crypto challenges (Base64/Base32/Caesar/Custom/Multi-layer/AES/HMAC/SM4/MD5 Signature)
- 🏢 **Realistic Business Scenario**: Multi-user VPS platform with company/individual accounts
- 🔐 **Encoding/Crypto Challenge System**: 9 encoding/crypto-related privilege escalation challenges (Base64/Base32/Caesar/Custom/Multi-layer/AES/HMAC/SM4/MD5 Signature), with automatic request parameter encoding in frontend
- 🔄 **Security/Vulnerable Mode Toggle**: Switch between modes for comparison learning
- 💡 **Progressive Hint System**: 3 levels of hints per vulnerability
- 📋 **Complete WriteUps**: Detailed vulnerability analysis and remediation
- 🔧 **One-Click Database Reset**: Restore initial state anytime
- 📦 **Single Binary Deployment**: Go compiles with embedded frontend assets
- 🚀 **Cross-Platform One-Click Scripts**: Windows (.bat/.ps1), Linux/macOS (.sh) double-click to run
- 🐳 **Docker Support**: One-command container deployment
- 🌙 **Theme Switching**: Light / Dark / System modes

## Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.25+ / Gin / GORM / SQLite (pure Go) / JWT |
| Frontend | Vue 3 / PrimeVue 4 / Pinia / Vite / TailwindCSS 4 |
| Database | SQLite3 (embedded) |
| Cryptography | golang.org/x/crypto / emmansun/gmsm (SM2/SM3/SM4) |

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
# Build
make build

# Run
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

## Directory Structure

```
OverstepLab/
├── src/
│   ├── backend/          # Go backend
│   │   ├── cmd/server/     # Entry point
│   │   ├── internal/       # Business code (handler/service/repository/model)
│   │   │   ├── crypto/     # Cryptography module (classical/modern/SM/signing)
│   │   │   ├── middleware/ # Middleware (JWT/RBAC/CORS/audit/encoding)
│   │   │   └── vuln/       # Vulnerability mode & challenge metadata
│   │   ├── database/       # Database layer (migration/seed)
│   │   └── router/         # Router
│   └── frontend/           # Vue3 frontend
│       └── src/
│           ├── api/        # API wrappers
│           ├── stores/     # Pinia state management
│           ├── views/      # Page views
│           └── router/     # Route configuration
├── doc/
│   ├── PRD.md              # Product Requirements Document
│   └── tutorials/          # Tutorials (8 articles)
├── README.md               # Chinese documentation
├── README_EN.md            # English documentation (this file)
├── SECURITY.md             # Security disclaimer
├── LICENSE                 # MIT License
├── AGENTS.md               # AI Agent Guide
├── start.bat               # Windows one-click startup (CMD)
├── start.ps1               # Windows one-click startup (PowerShell)
├── start.sh                # Linux/macOS one-click startup
├── Makefile
├── Dockerfile
└── docker-compose.yml
```

## Vulnerability List

### Traditional Privilege Escalation (13)

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

### Encoding / Crypto Challenges (9)

| ID | Category | Scenario | Difficulty |
|----|----------|----------|------------|
| E-01 | Encoding | Base64 encoded parameter bypass | ★☆☆ |
| E-02 | Encoding | Base32 encoded parameter bypass | ★☆☆ |
| E-03 | Encoding | Caesar cipher parameter offset bypass | ★★☆ |
| E-04 | Encoding | Custom Base64 alphabet bypass | ★★☆ |
| E-05 | Encoding | Multi-layer nested encoding bypass | ★★★ |
| E-06 | Encryption | AES-256-GCM encrypted parameter bypass | ★★★ |
| E-07 | Signature | HMAC-SHA256 signature verification bypass | ★★★★ |
| E-08 | Encryption | SM4-CBC national crypto parameter bypass | ★★★★ |
| E-09 | Signature | MD5 hash signature verification bypass | ★★★ |

## Crypto Tool Library

The frontend `src/utils/crypto.ts` contains built-in cryptographic utility functions, primarily used for automatic request encoding in encoding challenges, and also exposed to the browser console for manual invocation:

| Category | Supported Algorithms |
|----------|---------------------|
| Encoding | Base64, Base32, Base58, Hex |
| Classical | Caesar Cipher (ROT3) |
| Symmetric | AES-256-GCM, SM4-CBC (simplified implementation) |
| Hash/Signing | MD5, HMAC-SHA256 |
| Custom Encoding | Custom Base64 alphabet, Multi-layer nested encoding |

These utility functions can be manually invoked via `window.CryptoUtils` in the browser console.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `OVERSTEPLAB_PORT` | `5000` | Service listening port |
| `OVERSTEPLAB_DB_PATH` | `./oversteplab.db` | SQLite database path |
| `OVERSTEPLAB_SAFE_MODE` | `false` | Safe mode toggle (`true` enables all vulnerability patches) |
| `OVERSTEPLAB_JWT_SECRET` | Randomly generated | JWT signing secret |

## Contributing

Issues and Pull Requests are welcome. When adding new vulnerability scenarios, please update the PRD, tutorials, and WriteUp accordingly.

## License

MIT License. See [LICENSE](LICENSE).

## Security

This project contains **intentionally implanted vulnerabilities** for educational purposes only. See [SECURITY.md](SECURITY.md).
