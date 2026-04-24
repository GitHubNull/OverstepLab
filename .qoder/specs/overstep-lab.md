# OverstepLab Implementation Plan

## Context

OverstepLab is a deliberately vulnerable privilege escalation testing range (靶场) simulating a multi-tenant VPS management platform called "CloudNest". The project needs to be built from scratch according to the PRD in `doc/PRD.md`. It targets security researchers, penetration testers, developers, and trainers who need a realistic environment to practice horizontal IDOR, vertical escalation, and context-based access control bypasses.

**Key user requirements:**
- Backend and frontend code under `src/`
- Go binary embeds and serves frontend static assets (no separate nginx)
- Chinese-primary README + English README
- LICENSE, legal/security disclaimer, AGENTS.md
- Tutorial documents in `doc/tutorials/` with progressive difficulty

**Confirmed technical decisions:**
- SQLite driver: `modernc.org/sqlite` (pure Go, no CGO)
- Docker support: yes (Dockerfile + docker-compose.yml)
- Delivery: complete all milestones at once

---

## Project Structure

```
OverstepLab/
├── .gitignore
├── Makefile
├── Dockerfile
├── docker-compose.yml
├── LICENSE
├── README.md
├── README_EN.md
├── SECURITY.md
├── AGENTS.md
├── doc/
│   ├── PRD.md
│   └── tutorials/
│       ├── index.md
│       ├── 01-quickstart.md
│       ├── 02-horizontal-idor.md
│       ├── 03-vertical-escalation.md
│       ├── 04-context-escalation.md
│       ├── 05-advanced-combo.md
│       └── 06-writeups.md
├── src/
│   ├── backend/
│   │   ├── cmd/server/main.go
│   │   ├── internal/
│   │   │   ├── config/
│   │   │   ├── middleware/
│   │   │   ├── handler/
│   │   │   ├── service/
│   │   │   ├── repository/
│   │   │   ├── model/
│   │   │   ├── dto/
│   │   │   ├── common/
│   │   │   ├── vuln/
│   │   │   └── web/
│   │   │       ├── embed.go
│   │   │       ├── static.go
│   │   │       └── dist/          # frontend build output
│   │   ├── database/
│   │   │   ├── database.go
│   │   │   ├── migration/
│   │   │   └── seed/
│   │   ├── router/
│   │   └── go.mod
│   └── frontend/
│       ├── public/
│       ├── src/
│       │   ├── api/
│       │   ├── assets/
│       │   ├── components/
│       │   ├── composables/
│       │   ├── layouts/
│       │   ├── router/
│       │   ├── stores/
│       │   ├── types/
│       │   ├── views/
│       │   ├── App.vue
│       │   └── main.ts
│       ├── index.html
│       ├── package.json
│       ├── vite.config.ts
│       └── tsconfig.json
└── tmp/
```

---

## Implementation Phases

### Phase 1: Project Scaffolding & Infrastructure

**Goal:** Initialize both backend and frontend projects, establish build pipeline, and verify end-to-end connectivity.

1. **Backend scaffolding**
   - Create `src/backend/go.mod` with module path `github.com/oversteplab/oversteplab`
   - Add dependencies: `gin`, `gorm`, `modernc.org/sqlite`, `golang-jwt/jwt/v5`, `golang.org/x/crypto/bcrypt`, `swaggo` packages, `google/uuid`
   - Create `cmd/server/main.go` entry point
   - Create `internal/config/config.go` for env-based config (port, db path, JWT secret, safe mode)
   - Create `internal/common/` with response wrapper, error codes, utilities (bcrypt hash, random strings)

2. **Frontend scaffolding**
   - Initialize `src/frontend/` with Vite + Vue + TypeScript template via pnpm
   - Install dependencies: `vue@^3.4`, `vue-router@^4`, `pinia@^2`, `primevue@^4`, `@primevue/themes`, `primeicons`, `primeflex`, `axios`, `xterm`, `xterm-addon-fit`
   - Configure `vite.config.ts` with path alias `@`, proxy `/api` to `localhost:8080`, `build.outDir: '../../backend/internal/web/dist'`
   - Configure `tsconfig.json` with strict mode
   - Set up `main.ts` with PrimeVue, Pinia, Router
   - Create `App.vue` with global Toast and ConfirmDialog providers

3. **Build pipeline & root files**
   - Create `Makefile` with targets: `frontend`, `backend`, `build`, `dev-frontend`, `dev-backend`, `docker`, `clean`
   - Create `.gitignore`
   - Create `Dockerfile` (multi-stage: node builder -> go builder -> alpine runtime)
   - Create `docker-compose.yml`
   - Verify: `make build` produces a single binary that serves the embedded frontend

### Phase 2: Database Layer & Seed Data

**Goal:** Define all models, auto-migrate, and populate seed data on first startup.

1. **Models** (`src/backend/internal/model/`)
   - `user.go`: User struct with user_type, role, company_id, status
   - `company.go`: Company struct
   - `vps.go`: VPSInstance struct with owner_id, company_id, specs, status
   - `order.go`: Order struct with order_no, type, amount, status
   - `ticket.go`: Ticket struct
   - `ticket_reply.go`: TicketReply struct
   - `api_key.go`: APIKey struct with key_value (hash), key_prefix, permissions JSON, status, expire_at
   - `audit_log.go`: AuditLog struct with action, resource_type, resource_id, detail JSON, ip_address
   - `bill.go`: Bill struct with type, amount, balance_after

2. **Database init** (`src/backend/database/`)
   - `database.go`: SQLite connection init using `modernc.org/sqlite`, connection pooling
   - `migration/migration.go`: `AutoMigrate` all models
   - `seed/seed.go`: Pre-seed in transaction: 2 companies, 9 users, 7 VPS instances, sample orders/bills/tickets
   - Seed data uses deterministic IDs (1-9 users, 1-2 companies) for reproducible exploits
   - `ResetDatabase()` function: close DB, delete file, re-init, migrate, seed

3. **Verification**
   - Run backend, confirm DB file created, confirm seed data present
   - Query users table to verify 9 pre-seeded accounts

### Phase 3: Authentication & Core Middleware

**Goal:** JWT auth, RBAC, audit logging, CORS, and vulnerability mode toggle.

1. **Middleware** (`src/backend/internal/middleware/`)
   - `auth.go`: Extract `Authorization: Bearer <token>`, validate JWT, inject user into gin context. **C-03 hook**: in vulnerable mode, skip API key status/expiration checks
   - `rbac.go`: Role-permission matrix. `RequirePermission(perm)` middleware. In secure mode enforce; in vulnerable mode allow bypass at service layer
   - `audit.go`: Post-handler audit log capture via `c.Next()` + async channel writer
   - `cors.go`: CORS headers for dev/prod
   - `recovery.go`: Panic recovery with unified error response

2. **Vulnerability mode** (`src/backend/internal/vuln/`)
   - `mode.go`: `atomic.Int32` global state (0=vulnerable default, 1=secure). `SetSecureMode()`, `IsSecureMode()`
   - Env var `OVERSTEPLAB_SAFE_MODE` initializes state on startup
   - Runtime toggle via admin API

3. **Auth handlers** (`src/backend/internal/handler/auth.go`)
   - `POST /api/v1/auth/register` (company + individual)
   - `POST /api/v1/auth/login` (returns JWT + refresh token)
   - `POST /api/v1/auth/refresh`
   - `POST /api/v1/auth/logout`

4. **Auth service & repository**
   - Password hashing with bcrypt
   - JWT generation/validation with `golang-jwt/jwt/v5`

5. **Verification**
   - Register a new user, login, access protected endpoint with token, verify 401 without token

### Phase 4: User, Company & VPS Modules

**Goal:** Implement the most vulnerability-dense modules.

1. **User module** (`handler/user.go`, `service/user.go`, `repository/user.go`)
   - `GET /api/v1/user/profile`, `PUT /api/v1/user/profile`, `PUT /api/v1/user/password`
   - `GET /api/v1/users/:id` — **H-02 target**: in vulnerable mode, no ownership check

2. **Company module** (`handler/company.go`, `service/company.go`)
   - `GET /api/v1/company/members`
   - `POST /api/v1/company/members` — **V-02, C-02 targets**: missing role/type checks in vulnerable mode
   - `PUT /api/v1/company/members/:id`
   - `DELETE /api/v1/company/members/:id`
   - `PUT /api/v1/company/members/:id/role` — **V-05 target**: self-role escalation possible

3. **VPS module** (`handler/vps.go`, `service/vps.go`, `repository/vps.go`)
   - `GET /api/v1/vps`, `POST /api/v1/vps`, `GET /api/v1/vps/:id` — **H-01, C-01**
   - `POST /api/v1/vps/:id/start/stop/restart` — **V-01**
   - `POST /api/v1/vps/:id/reinstall` — **V-03**
   - `DELETE /api/v1/vps/:id`
   - `GET /api/v1/vps/:id/console` (mock, returns placeholder terminal data)
   - Service layer branches on `vuln.IsSecureMode()` for each vulnerability

4. **Frontend: Auth & VPS pages**
   - `LoginView.vue`, `RegisterView.vue` with company/individual tabs
   - `DashboardView.vue` with stat cards
   - `VpsListView.vue`, `VpsDetailView.vue`, `VpsConsoleView.vue` (xterm.js mock terminal)
   - Pinia stores: `auth.ts`, `vps.ts`, `user.ts`
   - API layer: `api/client.ts`, `api/auth.ts`, `api/vps.ts`, `api/user.ts`
   - Router with guards, `MainLayout.vue` with sidebar and topbar

5. **Verification**
   - Login as `acme_ops`, list VPS, view detail, verify H-01 by changing ID parameter
   - Login as `acme_viewer`, verify V-01 by calling start/stop API directly
   - Toggle safe mode, verify same requests are blocked

### Phase 5: Finance, Tickets, API Keys & Audit Logs

**Goal:** Complete remaining business modules with their respective vulnerabilities.

1. **Finance module** (`handler/order.go`, `handler/bill.go`, `service/`, `repository/`)
   - `GET /api/v1/orders`, `GET /api/v1/orders/:id` — **H-03**
   - `GET /api/v1/bills`, `POST /api/v1/bills/recharge`, `GET /api/v1/bills/export` (CSV)

2. **Ticket module** (`handler/ticket.go`, `service/ticket.go`)
   - `GET /api/v1/tickets`, `POST /api/v1/tickets`, `GET /api/v1/tickets/:id` — **H-04**
   - `POST /api/v1/tickets/:id/reply` — **H-04**
   - `PUT /api/v1/tickets/:id/close`

3. **API Key module** (`handler/apikey.go`, `service/apikey.go`)
   - `GET /api/v1/apikeys`, `POST /api/v1/apikeys`, `DELETE /api/v1/apikeys/:id` — **H-05**

4. **Audit Log module** (`handler/auditlog.go`)
   - `GET /api/v1/audit-logs` (personal/company/platform based on role)

5. **Frontend pages**
   - `OrdersView.vue`, `BillsView.vue`, `RechargeView.vue`
   - `TicketListView.vue`, `TicketDetailView.vue`, `TicketCreateView.vue`
   - `ApiKeysView.vue`
   - `AuditLogsView.vue`
   - Corresponding Pinia stores and API modules

6. **Verification**
   - Test H-03, H-04, H-05 by modifying IDs in requests
   - Verify CSV export works
   - Verify audit logs capture actions

### Phase 6: Platform Admin & Challenge Panel

**Goal:** Admin endpoints, vulnerability challenge UI, hints, writeups, and reset functionality.

1. **Admin module** (`handler/admin.go`, `service/admin.go`)
   - `GET /api/v1/admin/users`, `PUT /api/v1/admin/users/:id/status` — **V-04**
   - `GET /api/v1/admin/companies`
   - `GET /api/v1/admin/vps`
   - `GET /api/v1/admin/audit-logs`
   - `POST /api/v1/admin/reset` — database reset

2. **Challenge module** (`handler/challenge.go`)
   - `GET /api/v1/challenges` — list all challenges with completion status
   - `GET /api/v1/challenges/:id` — detail
   - `GET /api/v1/challenges/:id/hints/:level` — progressive hints
   - `POST /api/v1/challenges/:id/complete` — mark completed
   - `GET /api/v1/security-mode`, `PUT /api/v1/security-mode`

3. **Challenge metadata** (`internal/vuln/challenges.go`)
   - Define all 13 challenges: ID, title, category, difficulty, description, 3 hints each, writeup text

4. **Frontend: Admin & Challenges**
   - `AdminUsersView.vue`, `AdminCompaniesView.vue`, `AdminVpsView.vue`, `AdminSystemView.vue` (with reset button)
   - `ChallengesView.vue` — challenge cards, hint reveal, writeup display, completion tracking
   - Security mode toggle switch in `AppTopbar.vue`
   - Pinia stores: `admin.ts`, `challenges.ts`

5. **Verification**
   - Login as `admin`, verify admin panel shows all users/companies/VPS
   - Login as `alice`, verify V-04 by calling admin APIs directly
   - Open challenge panel, verify all 13 challenges listed, hints reveal progressively
   - Toggle security mode, verify behavior changes
   - Click reset, verify database restores to seed state

### Phase 7: Documentation

**Goal:** Create all required project documentation.

1. **Root documentation**
   - `README.md`: security warning banner, project intro, features, tech stack, quick start (dev + production), pre-seeded accounts table, directory structure, vulnerability quick reference, contribution guide, license link
   - `README_EN.md`: English mirror of README.md
   - `LICENSE`: MIT License
   - `SECURITY.md`: legal disclaimer (intentionally vulnerable, isolated use only, no production deployment, no liability), vulnerability list, deployment restrictions, responsible disclosure
   - `AGENTS.md`: project identity, architecture overview, directory conventions, key files reference, coding conventions (Go layered, Vue Composition API), intentional vulnerability protocol, safe mode explanation, development workflow

2. **Tutorial documents** (`doc/tutorials/`)
   - `index.md`: learning path overview, prerequisites, tools recommended, links to all tutorials
   - `01-quickstart.md`: install, build, login with 3 sample accounts, UI tour, role matrix
   - `02-horizontal-idor.md`: H-01~H-05 walkthroughs, browser Network tab, curl basics
   - `03-vertical-escalation.md`: V-01~V-05 walkthroughs, UI hiding vs backend enforcement
   - `04-context-escalation.md`: C-01~C-03 walkthroughs, business logic flaws, state validation
   - `05-advanced-combo.md`: chaining vulnerabilities, attack graphs, automation scripts
   - `06-writeups.md`: per-vulnerability detailed writeup (ID, category, difficulty, affected endpoints, root cause, exploitation steps, impact, remediation, safe mode behavior)

3. **Verification**
   - All markdown files render correctly
   - Tutorial links work
   - Security warning is prominent

### Phase 8: Polish, Build Verification & Docker

**Goal:** Final integration, responsive design pass, production build verification, Docker test.

1. **Frontend polish**
   - Theme switching (dark/light)
   - Mobile responsiveness (collapsible sidebar, stacked layouts)
   - Loading states, skeletons, empty states
   - Toast notifications across all async operations
   - Form validation feedback

2. **Build verification**
   - `make clean && make build` produces working single binary
   - Run binary, verify all pages load, API works, static assets served
   - Test database reset and re-seed

3. **Docker verification**
   - `docker-compose up --build` succeeds
   - Access app at `http://localhost:8080`
   - Verify SQLite persists in `./data/`

4. **Final checks**
   - Verify all 13 vulnerabilities are exploitable in vulnerable mode
   - Verify all 13 are blocked in secure mode
   - Verify frontend build outputs directly to `src/backend/internal/web/dist/`
   - Verify Go embed serves SPA with fallback to index.html

---

## Critical Files to Modify/Create

| File | Purpose |
|------|---------|
| `src/backend/cmd/server/main.go` | Entry point, DB init, seed, router, server start |
| `src/backend/internal/web/embed.go` | `//go:embed all:dist` directive |
| `src/backend/internal/web/static.go` | SPA static file serving with index.html fallback |
| `src/backend/internal/vuln/mode.go` | Security/vulnerable mode atomic state |
| `src/backend/internal/middleware/auth.go` | JWT validation, C-03 vulnerability hook |
| `src/backend/internal/middleware/rbac.go` | Role-based access control |
| `src/backend/database/seed/seed.go` | All pre-seeded test data |
| `src/backend/internal/service/vps.go` | H-01, V-01, V-03, C-01 vulnerability branching |
| `src/backend/internal/service/company.go` | V-02, V-05, C-02 vulnerability branching |
| `src/backend/internal/service/user.go` | H-02 vulnerability branching |
| `src/frontend/vite.config.ts` | Build output to backend embed dir, API proxy |
| `src/frontend/src/api/client.ts` | Axios instance with JWT interceptors |
| `src/frontend/src/router/guards.ts` | Auth and role-based navigation guards |
| `src/frontend/src/stores/auth.ts` | Auth state, security mode toggle |
| `src/frontend/src/layouts/MainLayout.vue` | App shell with sidebar, topbar, security switch |
| `Makefile` | Orchestrates zero-copy build pipeline |
| `Dockerfile` | Multi-stage container build |
| `README.md` / `README_EN.md` | Project documentation |
| `SECURITY.md` | Legal/security disclaimer |
| `AGENTS.md` | AI agent guidance |
| `doc/tutorials/*.md` | Progressive learning tutorials |

---

## Verification Plan

1. **Backend unit tests**
   - Run `go test ./...` in `src/backend/` after implementing each module
   - Test auth flow: register -> login -> access protected -> logout -> 401
   - Test each vulnerability in both modes using HTTP client

2. **Frontend type check**
   - Run `pnpm type-check` in `src/frontend/`
   - Ensure no TypeScript errors

3. **End-to-end manual verification**
   - Build single binary: `make build`
   - Run: `./bin/oversteplab` (or equivalent on Windows)
   - Open browser to `http://localhost:8080`
   - Login with each pre-seeded account, verify role-appropriate UI
   - Test vulnerability exploitation for H-01~H-05, V-01~V-05, C-01~C-03
   - Toggle safe mode, verify exploits are blocked
   - Test database reset
   - Test Docker deployment: `docker-compose up --build`

4. **Documentation verification**
   - Preview all markdown files
   - Confirm tutorial index links are valid
   - Confirm security warning is present in README and SECURITY.md
