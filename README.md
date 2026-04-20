<div align="center">

# ZATRANO

[![GitHub — zatrano/zatrano](https://img.shields.io/badge/GitHub-zatrano%2Fzatrano-181717?style=for-the-badge&logo=github)](https://github.com/zatrano/zatrano)
[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/dl/)
[![Fiber](https://img.shields.io/badge/Fiber-v3-00ACFF?style=for-the-badge)](https://github.com/gofiber/fiber)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![MySQL](https://img.shields.io/badge/MySQL-4479A1?style=for-the-badge&logo=mysql&logoColor=white)](https://www.mysql.com/)
[![SQLite](https://img.shields.io/badge/SQLite-003B57?style=for-the-badge&logo=sqlite&logoColor=white)](https://www.sqlite.org/)
[![SQL Server](https://img.shields.io/badge/SQL%20Server-CC2927?style=for-the-badge&logo=microsoftsqlserver&logoColor=white)](https://www.microsoft.com/sql-server)
[![Redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io/)
<br>
[![GORM](https://img.shields.io/badge/GORM-632CA6?style=for-the-badge)](https://gorm.io/)
[![Zap](https://img.shields.io/badge/Zap-structured%20logs-121212?style=for-the-badge)](https://github.com/uber-go/zap)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-6BA539?style=for-the-badge&logo=openapiinitiative&logoColor=white)](https://www.openapis.org/)
[![GraphQL](https://img.shields.io/badge/GraphQL-E10098?style=for-the-badge&logo=graphql&logoColor=white)](https://graphql.org/)
[![gqlgen](https://img.shields.io/badge/gqlgen-311C87?style=for-the-badge)](https://github.com/99designs/gqlgen)
[![golang-migrate](https://img.shields.io/badge/golang--migrate-SQL-00599C?style=for-the-badge)](https://github.com/golang-migrate/migrate)
[![Cobra CLI](https://img.shields.io/badge/Cobra-CLI-7E43B6?style=for-the-badge)](https://github.com/spf13/cobra)
[![Viper](https://img.shields.io/badge/Viper-config-273F5B?style=for-the-badge)](https://github.com/spf13/viper)
[![AWS SDK](https://img.shields.io/badge/AWS-S3%20SDK-232F3E?style=for-the-badge&logo=amazonaws&logoColor=white)](https://aws.amazon.com/sdk-for-go/)
[![OAuth2](https://img.shields.io/badge/OAuth2-x--oauth2-4285F4?style=for-the-badge)](https://pkg.go.dev/golang.org/x/oauth2)

</div>

---

## Overview

**ZATRANO** is an opinionated **Go backend platform** for teams that want **one consistent way** to ship, configure, migrate, and operate **modular monoliths** or **HTTP-first services**—without rebuilding sessions, queues, mail plumbing, and OpenAPI docs in every repository.

### What ZATRANO is

The project combines a **runtime** (Fiber), **persistence** (GORM across multiple SQL engines with **golang-migrate**), **Redis** patterns for **sessions**, **cache**, and **queues**, **mail & notifications**, an **event bus**, **OpenAPI** documentation, **optional GraphQL**, **multi-tenancy**, **audit**, **feature flags**, **real-time broadcasting** (WebSocket / SSE), and production-minded defaults—together with a **CLI** (`zatrano`) and **generators** that scaffold modules, CRUD, views, jobs, policies, and **wire** markers so structure stays uniform across teams.

Application code imports the **public contract** from **`pkg/`**. This repository also ships **`internal/cli`**, **`internal/gen`**, and **embedded SQL migrations** consumed by **`zatrano db migrate`**.

### What it is not

| | |
|---|---|
| **Not a minimal router** | You adopt **defaults and conventions** (sessions/CSRF, RBAC building blocks, migration layout, OpenAPI + Scalar, and more). The goal is **predictable delivery**, not a blank slate. |
| **Not your domain model** | You receive **layering** (handler → service → repository) and cross-cutting helpers; **business rules**, aggregates, and external integrations remain **your code**. |
| **Not a microservices framework** | The sweet spot is **modular monoliths** and **well-factored HTTP services**. It does **not** prescribe service mesh, orchestration, or a particular deployment topology. |

### Adoption path

1. **Add the dependency** — require `github.com/zatrano/zatrano`, bootstrap with **`pkg/config`** and **`pkg/core`**, register routes through **`zatrano.Start`** (or use **`zatrano.Run()`** when you have no custom routes).
2. **Scaffold and wire** — use **`zatrano new`** / **`zatrano gen …`** plus **`zatrano:wire`** markers so modules mount from a single registration point.
3. **Run and observe** — apply schema with **`zatrano db migrate`**, validate configuration with **`doctor`** / **`config validate`**, expose **`/health`** / **`/ready`**, and use the sections below for queues, tenants, mail, and backups.

### Stack at a glance

| | |
|---|---|
| **Module** | `github.com/zatrano/zatrano` |
| **Go** | 1.25+ |
| **HTTP & APIs** | Fiber v3, OpenAPI 3, optional GraphQL (gqlgen) |
| **Data** | GORM + supported SQL engines ([Database](#database)); schema changes via **golang-migrate** (`zatrano db migrate` / `rollback`) |
| **Shared infrastructure** | Redis (sessions, cache, queues), Zap structured logging, optional AWS S3 SDK, OAuth2 (`x/oauth2`) |

> **Status:** Active development. Treat **`pkg/`** as the **stable public surface** for applications that embed the platform.

### Developer & Project Owner:

**Serhan KARAKOÇ** — [github.com/serhankarakoc](https://github.com/serhankarakoc)

---

## Table of Contents

- [Overview](#overview)
- [Features](#features-roadmap)
- [Feature catalog](#feature-catalog)
- [Layout](#layout-pkg-vs-internal)
- [Requirements](#requirements)
- [Database](#database)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [CLI Commands](#cli-commands)
- [HTTP Routes & Internationalization (i18n)](#http-current)
- [Validation](#validation)
- [Authorization (RBAC & Gate/Policy)](#authorization-rbac--gatepolicy)
- [Cache System](#cache-system)
- [Queue / Job System](#queue--job-system)
- [Mail System](#mail-system)
- [Event / Listener System](#event--listener-system)
- [Broadcasting / WebSocket](#broadcasting--websocket)
- [Multi-tenancy](#multi-tenancy)
- [Audit / Activity log](#audit--activity-log)
- [Full-text search](#full-text-search)
- [Feature flags](#feature-flags)
- [GraphQL](#graphql)
- [Scheduled tasks](#scheduled-tasks)
- [Admin dashboard](#admin-dashboard)
- [REST API helpers](#rest-api-helpers)
- [Outbound HTTP client](#outbound-http-client)
- [Notifications](#notifications)
- [HTTP testing helpers](#http-testing-helpers)
- [Repository / Data](#repository--data)
- [Storage / File Management](#storage--file-management)
- [View / Template System](#view--template-system)
- [Configuration](#configuration)
  - [Database migrations (SQL)](#database-migrations-sql)
- [Development](#development)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)

---

## Features (roadmap)

| Area | Plan |
|------|------|
| Architecture | Modular core + pluggable modules (modular monolith) |
| Layers | Handler → Service → Repository (mandatory bases) |
| Web | Fiber HTML templates, CSRF, **validation** (`go-playground/validator`), flash, **CORS**, **rate limit**, **i18n** (JSON locales), **cache** (Memory/Redis), security headers, gzip, static |
| **View Engine** | Layout inheritance (`{{extends}}`), block/section system, component partials, form builder helpers, flash messages, old-input repopulation, versioned asset URLs, Vite/esbuild manifest integration, HMR dev server proxy |
| API | REST + **OpenAPI 3** (`api/openapi.yaml`, `/docs`, `/openapi.yaml`), **Resource/Transformer** (model→JSON, hide sensitive fields, shape relations), **Standard response envelope** ({data, meta, links}, JSON:API compatible), **Cursor pagination** (keyset for large datasets), **Throttle** (user/JWT subject rate limiting, Redis counters), **API key management** (api_keys table, middleware, rotation), **Versioning manager** (v1/, v2/ auto groups, config-driven prefixes) |
| Auth | **Session (Redis) + CSRF**; **JWT** for `/api/v1/private/*`; **OAuth2** (Google/GitHub) browser login; **RBAC** (role→permission, DB-backed); **Gate/Policy** (resource-based authorization); **Password Reset** / **Email Verification** (transactional e-mail via **`pkg/notifications`** → **mail** channel on `App.Notifications`, backed by `App.Mail`); **Brute Force Protection** (IP+username rate limiting, Redis); **TOTP 2FA** (Google Authenticator compatible, QR code generation); **Session Management** (list/revoke active sessions, device info); **JWT Refresh Tokens** (token rotation, refresh token table) |
| Cache | **Memory / Redis** drivers, **Tag-based** invalidation, **Middleware** support |
| Queue | **Redis-backed** job queue, delayed jobs (ZADD), auto retry + exponential backoff, failed jobs (PostgreSQL) |
| **Scheduled Tasks** | Cron scheduling with `robfig/cron/v3`, fluent `schedule.Call(fn).Daily().At("08:00")`, `EveryMinute`/`Hourly`/`Daily`/`Weekly`/`Monthly`, Redis overlap lock |
| Mail | **SMTP / Log** drivers, HTML templates with layouts, queue integration, attachments, Mailable pattern |
| Events | **Sync and async** event bus, `ShouldQueue` for queue-backed listeners, `gen event` + `gen listener` |
| **Notifications** | Multi-channel delivery (Database, Mail, SMS, Push), read/unread tracking, **Twilio / Netgsm drivers**, **FCM / APNs**, `gen notification` |
| **Broadcasting** | **WebSocket hub** (`pkg/broadcast`, channel fan-out, `github.com/fasthttp/websocket` + Fiber v3), **private / presence** channels (JWT `sub`), **online list** (`Hub.OnlineOn`), **SSE** one-way push, **Pusher-compatible** wire format (Echo / pusher-js friendly) |
| **Multi-tenancy** | **ResolveTenant** middleware (header `X-Tenant-ID` or subdomain), **`tenant.FromContext` / Locals**, **row isolation** via `repository.NewTenantAware` + `TenantScope`, optional **`TenantFK`** embed; **schema isolation** via `tenant.GormSession` + **`zatrano db tenants`** (migrate/rollback/create-schema with PostgreSQL `search_path`) |
| **Audit** | **Model activity** (`zatrano_activity_logs`, GORM hooks, `audit.RegisterSubject`, JSON Patch **diff** via `audit.DiffJSONPatch`), **HTTP audit** (`middleware.AuditLog`, `zatrano_http_audit_logs` or JSONL **file**), **`audit.WithUser` / `WithRequest`** on `context.Context` |
| **Full-text search** | **PostgreSQL** `tsvector` / `plainto_tsquery` via **`repository.Scope`** helpers in **`pkg/search`**, **Meilisearch / Typesense** lightweight HTTP drivers, bulk **`zatrano search import <Model>`** with **`search.RegisterImporter`** |
| **Feature flags** | **`pkg/features`** — YAML and/or **`zatrano_feature_flags`** table, user + role + **percentage rollout** (A/B), **`app.Features.For(user).IsEnabled`**, **`middleware.RequireFeature`**, template **`{{if feature . "key"}}`** (with **`ViewData`**) |
| **GraphQL** | **gqlgen** schema-first (`api/graphql/*.graphqls`, `gqlgen.yml`), **`/graphql`** via Fiber **`adaptor`**, optional **GraphiQL** playground, **`graph-gophers/dataloader`** hooks (**`Loaders`**, **`WithLoaders`**), **`zatrano gen graphql <Model>`** |
| **Testing** | **HTTP test client** (Fiber.Test() wrapper, Get/Post/WithToken, AssertStatus/AssertJSON), **Database factory** (gofakeit-based test data, gen factory), **Transaction rollback** (TestSuite struct, SetupTest/TeardownTest), **In-memory cache driver** (no Redis required), **Mail fake** (captures emails in memory, assert sent), **Queue fake** (captures dispatched jobs, assert dispatched) |
| Data | **Generic Repository** pattern, automated soft-deletes, **chainable Scopes**, Offset-based pagination |
| DB / Ops | **PostgreSQL · MySQL · SQLite · SQL Server** (`database_driver` + GORM); **`zatrano db migrate` / `rollback`** (default **embed** SQL in **`pkg/migrations/sql/<driver>/`**), **`db seed`**, **`db backup` / `restore`** (Postgres CLI tools) |
| **Storage** | **Local / S3 / MinIO / Cloudflare R2** drivers, **signed URLs**, **image processing** (resize, crop, thumbnail), **Fiber middleware**, public + private disks |
| **HTTP Client** | Fluent JSON client with **WithToken**, **WithHeader**, **WithTimeout**, `Get`/`Post`/`Put`, automatic JSON marshal/unmarshal, retry on 5xx, and fake test transport |
| Ops | `/health`, `/ready`, `/status` |
| CLI | **`new`**, **`gen module`**, **`gen crud`**, **`gen admin`**, **`gen request`**, **`gen policy`**, **`gen job`**, **`gen mail`**, **`gen event`**, **`gen listener`**, **`gen notification`**, **`gen model`**, **`gen middleware`**, **`gen resource`**, **`gen test`**, **`gen seeder`**, **`gen factory`**, **`gen command`**, **`gen wire`**, **`gen view`**, **`gen graphql`**, `serve`, `db`, **`db tenants`**, **`schedule run`**, **`schedule list`**, **`search import`**, **`cache`**, **`queue`**, **`mail`**, **`openapi export`**, `openapi validate`, **`jwt sign`**, **`api-key create`**, **`api-key list`**, **`api-key revoke`**, … |

**Implemented now:** `serve`, `doctor`, **`routes`**, **`config print`**, **`config validate`**, **`verify`** (optional **`--race`**), `completion`, `version` / **`--version`**, **`new`**, **`gen module`** + **`gen crud`** + **`gen admin`** + **`gen request`** + **`gen policy`** + **`gen job`** + **`gen mail`** + **`gen event`** + **`gen listener`** + **`gen notification`** + **`gen model`** + **`gen middleware`** + **`gen resource`** + **`gen test`** + **`gen seeder`** + **`gen factory`** + **`gen command`** + **`gen wire`** + **`gen view`** + **`gen graphql`**, **`db`** (golang-migrate; default **embed** SQL from **`pkg/migrations/sql/<driver>/`**, optional **file** + `migrations_dir` / `--migrations`) + **`db tenants`** (per-tenant PostgreSQL schema migrate/rollback/create-schema), **`pkg/schedule`** + **`schedule run`** / **`schedule list`** (cron, Redis overlap lock), **`search import`** (Meilisearch/Typesense bulk index via `RegisterImporter`), **`pkg/features`** (flags, rollout, template + HTTP middleware), **`pkg/graphql`** (gqlgen + dataloader hooks), **`cache`** (Memory/Redis, Tags, middleware), **`queue`** (Redis FIFO, delayed jobs, retry, failed jobs, worker), **`mail`** (SMTP/log, templates, queue, attachments, preview), **`events`** (sync/async dispatch, ShouldQueue, queue-backed listeners), **`notifications`** (multi-channel, Database/SMS/Push, read-tracking, Twilio/Netgsm/FCM/APNs), **`broadcast`** (WebSocket hub, Pusher-style protocol, private/presence JWT channels, SSE), **`audit`** (model activity + HTTP audit, JSON Patch diffs), **`pkg/search`** (PostgreSQL FTS scopes + external drivers), **`openapi validate`** + **`openapi export`**, **`jwt sign`**, **`pkg/storage`** (local/S3/MinIO/R2, signed URLs, image processing), **`pkg/api`** (response envelope, cursor pagination, API keys, route versioning), **`pkg/http`** (fluent outbound JSON client + test fake), **`pkg/testing`** (Fiber HTTP test client, DB transaction test suite), **`pkg/admin`** (ops dashboard: metrics, logs tail/download when `admin.enabled`), **OAuth2**, **`http.*`** (CORS, rate limit, request timeout, body limit), **`i18n`** (JSON locales + Fiber helpers), **validation** (generic `Validate[T]`, i18n errors, custom rules, form requests), **authorization** (RBAC role→permission, Gate/Policy, `middleware.Can`, i18n 403), **multi-tenancy** (`middleware.ResolveTenant`, `pkg/tenant`, tenant-scoped repository), **view engine** (`{{extends}}` layout inheritance, `{{block}}` sections, `views/components/` partials, form builder, flash messages, old-input `{{old}}`, `{{asset}}` versioned URLs, Vite/esbuild manifest + HMR), Redis session + CSRF, JWT, Scalar **`/docs`**, **Air** (`.air.toml`).

## Feature catalog

Each item states **what it is for**, shows a **tiny config / CLI / Go** example, then points to the **full chapter** below.

### Database & migrations

**Purpose:** GORM connection + versioned SQL via **golang-migrate** (embed or on-disk tree).

```yaml
database_driver: postgres
database_url: postgres://user:pass@localhost:5432/app?sslmode=disable
migrations_source: embed # or file + migrations_dir
```

```bash
zatrano db migrate
zatrano db rollback --steps 1
```

**Details:** [Database](#database) · [Database migrations (SQL)](#database-migrations-sql)

### HTTP, OpenAPI & docs

**Purpose:** Fiber stack, merged **`/openapi.yaml`**, Scalar **`/docs`**, health routes, JWT/OAuth stubs, broadcast paths when enabled.

```bash
curl -s http://127.0.0.1:8080/api/v1/public/ping
zatrano openapi validate --merged
```

**Details:** [HTTP (current)](#http-current)

### Internationalization (i18n)

**Purpose:** JSON locale files under **`locales_dir`**; `i18n.T` / `i18n.Tf` in handlers; validation/auth messages can reuse the same bundle.

```go
msg := i18n.T(c, "app.welcome")
```

**Details:** [HTTP (current)](#http-current)

### Validation

**Purpose:** `go-playground/validator` tags + `validation.Validate[T]`; **`zatrano gen request`** scaffolds form requests.

**Details:** [Validation](#validation)

### Authorization (RBAC & Gate/Policy)

**Purpose:** DB-backed roles/permissions, resource policies, **`middleware.Can`** / role helpers; **`zatrano gen policy`**.

**Details:** [Authorization (RBAC & Gate/Policy)](#authorization-rbac--gatepolicy)

### Cache

**Purpose:** Memory or Redis-backed **`app.Cache`** (`Set`/`Get`, JSON helpers, `Remember`, tags, HTTP response caching).

```go
_ = app.Cache.Set(ctx, "stats:last", "42", 5*time.Minute)
```

**Details:** [Cache System](#cache-system)

### Queue & jobs

**Purpose:** Redis-backed FIFO, delayed jobs, retries; **`queue.Job`** + **`app.Queue.Dispatch`**; worker via **`zatrano queue work`**.

```go
_ = app.Queue.Dispatch(ctx, &jobs.SendWelcomeEmailJob{UserID: 1})
```

**Details:** [Queue / Job System](#queue--job-system)

### Mail

**Purpose:** SMTP/log mail with templates, queue integration, attachments; **`zatrano gen mail`**, **`zatrano mail preview`**.

**Details:** [Mail System](#mail-system)

### Notifications

**Purpose:** Multi-channel **`App.Notifications`** (`Register`, `Send`, `SendToChannels`); built-in mail/database/SMS/push channel types; **`zatrano gen notification`**.

```go
_ = app.Notifications.SendToChannels(ctx, notifications.NewNotification("Hi", "Body", "user@example.com"), "mail")
```

**Details:** [Notifications](#notifications) · (Mail channel wiring: [Mail System](#mail-system))

### Events

**Purpose:** **`app.Events.Listen` / `Fire`**; listeners may implement **`ShouldQueue`** for async replay via the queue.

```go
app.Events.Listen("user.created", events.ListenerFunc(func(ctx context.Context, e events.Event) error {
    return nil
}))
```

**Details:** [Event / Listener System](#event--listener-system)

### Broadcasting (WebSocket / SSE)

**Purpose:** In-process hub, Pusher-compatible WebSocket JSON, optional SSE; enable with **`broadcast`** in YAML.

```yaml
broadcast:
  enabled: true
```

**Details:** [Broadcasting / WebSocket](#broadcasting--websocket)

### Multi-tenancy

**Purpose:** Resolve tenant from header/subdomain, **`tenant.FromContext`**, tenant-aware repositories and optional per-tenant PostgreSQL schemas (**`zatrano db tenants`**).

**Details:** [Multi-tenancy](#multi-tenancy)

### Audit

**Purpose:** Model activity log + HTTP audit middleware; JSON Patch diffs; user/request on **`context.Context`**.

**Details:** [Audit / Activity log](#audit--activity-log)

### Full-text search

**Purpose:** PostgreSQL FTS **`repository.Scope`** helpers + Meilisearch/Typesense drivers; bulk **`zatrano search import`**.

**Details:** [Full-text search](#full-text-search)

### Feature flags

**Purpose:** YAML/DB flags, percentage rollout, **`middleware.RequireFeature`**, template **`{{if feature . "key"}}`**.

**Details:** [Feature flags](#feature-flags)

### GraphQL

**Purpose:** gqlgen schema under **`api/graphql/`**, Fiber adaptor at **`/graphql`**, **`zatrano gen graphql`**, dataloader hooks.

**Details:** [GraphQL](#graphql)

### Scheduled tasks (cron)

**Purpose:** Register cron jobs with **`pkg/schedule`** (`schedule.Call(fn).Daily().At("09:00")`, `WithoutOverlap()` + Redis); run **`zatrano schedule run`** or list with **`schedule list`**.

```go
_, _ = schedule.Call(func(ctx context.Context) error { return nil }).Name("heartbeat").EveryMinute().Register()
```

**Details:** [Scheduled tasks](#scheduled-tasks)

### Admin (ops UI + CRUD scaffold)

**Purpose:** Framework **`pkg/admin`** mounts dashboard, metrics, and log viewer under **`admin.path_prefix`** when **`admin.enabled`**; **`zatrano gen admin`** scaffolds module HTML list pages that reuse **`pkg/admin`** helpers.

**Details:** [Admin dashboard](#admin-dashboard)

### REST API helpers (`pkg/api`)

**Purpose:** Standard JSON **`api.Success` / `api.JSON`** envelope, cursor pagination helpers, **`api.NewVersionManager`**, hashed **`api_keys`** via **`api.NewKeyManager`**.

**Details:** [REST API helpers](#rest-api-helpers)

### Outbound HTTP client (`pkg/http`)

**Purpose:** Fluent **`http.NewClient()`** with **`WithToken`**, **`WithRetry`**, JSON **`Get`/`Post`/`Put`** helpers and **`http.Fake`** for tests.

**Details:** [Outbound HTTP client](#outbound-http-client)

### Repository & data access

**Purpose:** Generic **`repository.New[T]`**, chainable scopes, pagination helpers over GORM.

**Details:** [Repository / Data](#repository--data)

### Storage (files & disks)

**Purpose:** **`pkg/storage`** disks (local/S3/MinIO/R2), signed URLs, image helpers; symlink public assets with CLI.

```bash
zatrano storage:link
```

**Details:** [Storage / File Management](#storage--file-management)

### View / templates

**Purpose:** Layout **`{{extends}}`**, blocks, partials, form helpers, Vite manifest; **`zatrano gen view`**.

**Details:** [View / Template System](#view--template-system)

### HTTP integration testing (`pkg/testing`)

**Purpose:** **`testing.NewHTTPClient(app)`** for Fiber **`app.Test`** flows; **`testing.NewTestSuite(db)`** for transaction rollback around tests.

**Details:** [HTTP testing helpers](#http-testing-helpers)

### Tooling (CLI, config, quality gates)

**Purpose:** **`zatrano doctor`**, **`config validate`**, **`verify`**, OpenAPI export, JWT signing, **`gen wire`**, and the rest of **`zatrano gen …`**.

```bash
zatrano doctor
zatrano verify
```

**Details:** [CLI commands](#cli-commands) · [Configuration](#configuration) · [Development](#development)

---

## Layout (`pkg/` vs `internal/`)

| Path | Role |
|------|------|
| `pkg/config`, `pkg/core`, `pkg/server`, `pkg/health`, `pkg/middleware`, `pkg/security`, `pkg/auth`, `pkg/cache`, `pkg/queue`, `pkg/schedule`, `pkg/mail`, `pkg/notifications`, `pkg/events`, `pkg/broadcast`, `pkg/tenant`, `pkg/audit`, `pkg/search`, `pkg/features`, `pkg/graphql`, `pkg/oauth`, `pkg/openapi`, `pkg/i18n`, `pkg/validation`, `pkg/storage`, `pkg/api`, `pkg/http`, `pkg/admin`, `pkg/database`, `pkg/migrations` (embedded SQL; not a Go import target), `pkg/testing`, `pkg/zatrano`, `pkg/meta` | **Public API (`pkg/`)** — import from application code |
| `internal/cli`, `internal/db`, `internal/gen` | **CLI & code generators** — not imported by running services |

Scaffolded applications call **`zatrano.Start`** with **`RegisterRoutes: routes.Register`** (see `internal/routes/register.go`). Use **`zatrano.Run()`** when you do not inject custom routes.

---

## Requirements

- Go **1.25.0** or later
- **SQL database** — GORM and `zatrano db migrate`; see [Database](#database)
- **Redis** — sessions and CSRF (optional in local development; typically required in production when `redis_url` is set)
- **PostgreSQL client tools** — for `zatrano db backup` / `db restore` only: `pg_dump`, `pg_restore`, `psql` must be available on **`PATH`**

---

## Database

Application data access uses **GORM** (`pkg/database`). Schema revisions are applied or rolled back with **`zatrano db migrate`** / **`db rollback`** through **golang-migrate**.

### Supported engines

| Engine | Config (`database_driver`) | Notes |
|--------|---------------------------|--------|
| **PostgreSQL** | `postgres` (default when unset / empty) | Primary development target |
| **MySQL** | `mysql` | |
| **SQLite** | `sqlite` | Handy for local tools and tests |
| **SQL Server** | `sqlserver` | URL must use the `sqlserver://` form expected by `go-mssqldb` |

Set **`database_driver`** and **`database_url`** in YAML or environment variables. Examples live in **`config/examples/dev.yaml`**; validation and normalization are in **`pkg/config`**.

### Migration SQL layout

Versioned **`*.up.sql`** / **`*.down.sql`** for each engine live under **[`pkg/migrations/sql/`](pkg/migrations/sql/)** (`postgres/`, `mysql/`, `sqlite/`, `sqlserver/`). With **`migrations_source: embed`** (default), the CLI uses the embedded tree for the configured driver. Use **`migrations_source: file`** and **`migrations_dir`** for on-disk migrations (typical for scaffolded apps), or **`--migrations <dir>`** on **`db migrate`** / **`db rollback`** / **`db tenants …`** to point at a folder for that run only.

For **`migrations_source`**, seeds, and **`zatrano gen model`** paths, see **[Database migrations (SQL)](#database-migrations-sql)** under [Configuration](#configuration).

---

## Installation

Install the CLI globally:

```bash
go install github.com/zatrano/zatrano/cmd/zatrano@latest
```

---

## Quick start

Create a new app:

```bash
zatrano new app
cd app
zatrano serve
```

Or run the framework directly:

```bash
go run ./cmd/zatrano serve
```

Optional:

```bash
cp config/examples/dev.yaml config/dev.yaml
cp .env.example .env
```

After **`DATABASE_URL`** (and optional **`DATABASE_DRIVER`**; see [Database](#database)), apply the built-in schema (defaults to **embedded** migrations — `migrations_source: embed`):

```bash
zatrano db migrate --env dev --config-dir config
```

Validate or export OpenAPI (export merges `api/openapi.yaml` with framework routes — same as live `/openapi.yaml`):

```bash
go run ./cmd/zatrano openapi validate api/openapi.yaml
go run ./cmd/zatrano openapi validate --merged
go run ./cmd/zatrano openapi export --output api/openapi.merged.yaml
```

---

## CLI commands

| Command | Purpose |
|---------|---------|
| `zatrano serve` | HTTP server (`--addr`, `--env`, `--config-dir`, `--no-dotenv`) |
| `zatrano doctor` | Config (incl. **HTTP** middleware summary) + Postgres/Redis checks |
| `zatrano routes` | Print routes (same config as `serve`; `--json`, `--all`, **`--group`**) |
| `zatrano config print` | Effective config, **masked** secrets; **`--paths-only`** short summary (default **lines**; `json` / `yaml`) |
| `zatrano config validate` | Load + **validate** only (no DB/Redis); **`--quiet`** / **`-q`** for CI exit code only |
| `zatrano new <name>` | Scaffold app (`--module`, `--output`, `--replace-zatrano` for local dev) |
| `zatrano db migrate` | Apply embedded driver-specific SQL from `pkg/migrations/sql/<driver>/` by default (`migrations_source: embed`); use `file` + `migrations_dir` or `--migrations` for disk-based SQL |
| `zatrano db rollback` | Roll back (`--steps`) |
| `zatrano db seed` | Run `db/seeds/*.sql` in one transaction (no-op if no `.sql` files) |
| `zatrano db backup` | `pg_dump` → file/dir (`--format`: custom, plain, or directory; `--output` or default under `backups/`) |
| `zatrano db restore` | `pg_restore` / `psql` (**requires `--yes`**, optional `--clean`) |
| `zatrano db tenants migrate` | Apply migrations with PostgreSQL **`search_path`** scoped to `tenant_<key>` schema (`--tenant` required; same flags as `db migrate` for `--env`, `--config-dir`, `--migrations`, `--steps`) |
| `zatrano db tenants rollback` | Roll back tenant-schema migrations (`--tenant`, `--steps`, …) |
| `zatrano db tenants create-schema` | `CREATE SCHEMA IF NOT EXISTS` for the computed tenant schema name |
| `zatrano search import <model>` | Bulk-index a model registered with `search.RegisterImporter` into Meilisearch or Typesense (`search.enabled`, `database_url`; `--env`, `--config-dir`, `--no-dotenv`) |
| `zatrano gen module <name>` | Scaffold `modules/<name>/`; **wires** + **`go fmt`** on wire file (`--skip-wire`, `--module-root`, `--out`, `--dry-run`) |
| `zatrano gen crud <name>` | Add CRUD stubs + **form request structs** (`requests/`); **wires** `RegisterCRUD` + **`go fmt`** (same flags) |
| `zatrano gen request <name>` | Generate form request structs only (`modules/<name>/requests/create_*.go`, `update_*.go`) |
| `zatrano gen policy <name>` | Generate authorization policy stub (`modules/<name>/policies/<name>_policy.go`) implementing `auth.Policy` with CRUD methods |
| `zatrano gen job <name>` | Generate queue job stub (`modules/jobs/<name>.go`) implementing `queue.Job` with Handle, Retries, Timeout |
| `zatrano gen mail <name>` | Generate Mailable struct + HTML template (`modules/mails/<name>_mail.go` + `views/mails/<name>.html`) |
| `zatrano gen event <name>` | Generate event struct (`modules/events/<name>_event.go`) implementing `events.Event` |
| `zatrano gen notification <name>` | Generate a notification stub (`modules/notifications/<name>.go`) for multi-channel delivery |
| `zatrano gen listener <name>` | Generate listener (`modules/listeners/<name>_listener.go`); use `--queued` for async |
| `zatrano gen model <name>` | Generate a model scaffold under `pkg/repository/models/` and PostgreSQL migration stubs under `pkg/migrations/sql/postgres/` |
| `zatrano gen middleware <name>` | Generate a Fiber middleware stub under `pkg/middleware/` |
| `zatrano gen resource <name>` | Generate an API resource transformer stub under `pkg/resources/` |
| `zatrano gen test <name>` | Generate handler and service test stubs under `tests/` |
| `zatrano gen seeder <name>` | Generate a SQL seed file under `db/seeds/` |
| `zatrano gen factory <name>` | Generate a test data factory stub under `pkg/factory/` |
| `zatrano gen command <name>` | Generate a Cobra CLI command scaffold under `internal/cli/` |
| `zatrano gen wire <name>` | **Wire only** (no overwrite); picks `Register` / `RegisterCRUD` from existing files (`--register-only`, `--crud-only`) |
| `zatrano gen view <n>` | Scaffold server-rendered HTML templates under `views/<n>/` (`index.html`, `show.html`; `--with-form` adds `create.html` + `edit.html`; `--layout`, `--dry-run`) |
| `zatrano gen graphql <model>` | Add `api/graphql/<model>_stub.graphqls` + run **`go run github.com/99designs/gqlgen@v0.17.78 generate`** (`--module-root`, `--dry-run`, `--skip-generate`) |
| `zatrano openapi validate [path]` | Validate one file, or **`--merged`** (same as live `/openapi.yaml`; `--base`, optional positional overrides base) |
| `zatrano openapi export` | Write merged YAML (`--base`, `--output` or `-` for stdout) |
| `zatrano jwt sign` | Print HS256 token (`--sub`, `--secret`, config flags) |
| `zatrano cache clear` | Clear all cache or specific tags (`--tag`) |
| `zatrano queue work` | Start queue worker process (`--queue`, `--tries`, `--timeout`, `--sleep`) |
| `zatrano queue failed` | List failed jobs |
| `zatrano schedule run` | Start the scheduled task runner using registered tasks |
| `zatrano schedule list` | List registered scheduled tasks and cron expressions |
| `zatrano queue retry [id]` | Retry a failed job or `--all` |
| `zatrano queue flush` | Delete all failed jobs |
| `zatrano mail preview [name]` | Preview email template in browser (`--port`, `--layout`) |
| `zatrano storage:link` | Create symlink from `storage/app/public` to `public/storage` (`--force`, `--storage-path`, `--public-path`) |
| `zatrano storage:clear [disk]` | Clear all files from storage disk (`--force` to skip confirmation) |
| `zatrano completion …` | Shell completions |
| `zatrano verify` | **`go vet` + `go test` + merged OpenAPI** (PR/CI; `--race` for data races; `--no-vet`, `--no-test`, `--no-openapi`, `--module-root`) |
| `zatrano version` | Version string (also **`zatrano --version`**) |

**Windows / paths with spaces:** use `--replace-zatrano` pointing at your checkout; the scaffold quotes the path in `go.mod` when needed.

---

## HTTP (current)

| Method | Path | Notes |
|--------|------|-------|
| GET | `/` | JSON index (`env`, `endpoints`, `http` flags for CORS/rate limit, `error_includes_request_id`) |
| GET | `/health`, `/ready`, `/status` | Liveness / readiness / aggregate (`/status` includes `env`) |
| GET | `/openapi.yaml` | **Merged** OpenAPI (your file + built-in ops; **`/`** and **`/status`** include JSON schemas) |
| GET | `/docs` | Scalar API reference (CDN) |
| GET | `/api/v1/public/ping` | Public JSON |
| GET | `/api/v1/private/me` | **Bearer JWT** required if `jwt_secret` set |
| POST | `/api/v1/auth/token` | **Only if** `security.demo_token_endpoint: true` (blocked when `env: prod`) |
| GET | `/auth/oauth/google/login`, `/auth/oauth/github/login` | Starts OAuth2 (requires `oauth.enabled` + provider keys) |
| GET | `/auth/oauth/google/callback`, `/auth/oauth/github/callback` | OAuth redirect handler |
| GET | `/broadcast/ws` | **WebSocket** (when `broadcast.enabled: true`); Pusher-style JSON; JWT via query `access_token` or `Authorization` |
| GET | `/broadcast/sse/:channel` | **SSE** (when `broadcast.enabled` + `broadcast.sse_enabled`); same channel names as WebSocket; token via query or header |
| GET | `{admin.path_prefix}/` | **Admin** HTML dashboard (**`admin.enabled`**, **`View`** required; default prefix **`/admin`**) |
| GET | `{admin.path_prefix}/metrics` | Admin metrics (Redis queue depths from **`admin.queue_names`**, etc.) |
| GET | `{admin.path_prefix}/logs` | Tail of **`admin.log_file`** (if set) |
| GET | `{admin.path_prefix}/logs/download` | Download the configured log file |

**Session + CSRF:** enabled when `redis_url` is set and `security.session_enabled` / `csrf_enabled` are true. CSRF is skipped for `Authorization: Bearer …`, `csrf_skip_prefixes` (default includes `/api/`), and **`/auth/oauth/`** (OAuth callbacks).

**OAuth2:** set `oauth.enabled`, `oauth.base_url`, and `oauth.providers.google` / `github` client IDs. Redirect URLs in the provider console must be `{base_url}/auth/oauth/google/callback` (and the same for `github`). Session keys after login: `oauth_provider`, `oauth_subject`, `oauth_name`, `oauth_email`.

**Errors:** JSON responses use `{ "error": { "code", "message", "request_id"? } }`. `request_id` matches the **`X-Request-ID`** header when middleware runs (use it in logs and support tickets).

**HTTP middleware (`http` in YAML / `HTTP_*` env):**

- **CORS** — `http.cors_enabled`, `cors_allow_origins`, `cors_allow_methods`, `cors_allow_headers`, `cors_expose_headers`, `cors_allow_credentials`, `cors_max_age`. Default **off**. You cannot combine **`cors_allow_credentials: true`** with a wildcard origin **`*`** (browser rules); validation fails if you try.
- **Rate limit** — `http.rate_limit_enabled`, `rate_limit_max`, `rate_limit_window`, optional **`rate_limit_redis: true`** (uses **`redis_url`**; required if you enable Redis-backed limiting). Otherwise **in-memory** per process. Responses **under** the limit include **`X-RateLimit-*`** headers. When exceeded, **429** uses the same JSON `error` shape and Fiber sets **`Retry-After`** (RFC 6585).
- **Request timeout** — `http.request_timeout` (e.g. `60s`): Fiber **timeout** middleware; **408** JSON on overrun.
- **Body limit** — `http.body_limit` bytes (maps to Fiber **`BodyLimit`**; `0` = Fiber default **4 MiB**).
- **Graceful HTTP shutdown** — `http.shutdown_timeout` (default `15s`): upper bound for Fiber `ShutdownWithContext`. Use `zatrano.StartOptions.ShutdownHooks` for extra steps in the same deadline.
- **Zero-downtime restart (Unix)** — `http.graceful_restart: true`: [tableflip](https://github.com/cloudflare/tableflip) listener handoff; send **`SIGUSR2`** to trigger `Upgrade()`. Optional `http.graceful_restart_pid_file` for systemd-style setups. Requires a **real compiled binary** (not `go run`). On Windows the flag is ignored at runtime.

Order in the stack (see **`pkg/server/mount.go`**): **recover → request-id → (tenant, if enabled) → i18n (if enabled) → CORS → request timeout → rate limit → helmet → compress → session/CSRF → access log → (HTTP audit, if enabled) → view flash (if enabled) → feature locals (if enabled) → routes** (then modules, **`admin.Register`**, your **`RegisterRoutes`**, OAuth, OpenAPI, static).

**Access logging:** **`middleware.AccessLog`** writes one structured Zap line per request (method, path, status, latency, IP, **`request_id`**).

**Health:** **`health.Register`** exposes **`/health`**, **`/ready`**, **`/status`**; use **`health.Probe(ctx, app)`** for the same checks without HTTP.

---

## Broadcasting / WebSocket

ZATRANO ships an **in-memory broadcast hub** under **`pkg/broadcast`**: channel-based fan-out to **WebSocket** clients and optional **Server-Sent Events (SSE)** subscribers on the same channel names. The wire format follows a **Pusher-compatible subset** so frontends can reuse **Laravel Echo**, **pusher-js**, or any client that speaks `pusher:subscribe` / `pusher:connection_established` style JSON.

### Enable

```yaml
# config/dev.yaml
broadcast:
  enabled: true
  path_prefix: /broadcast          # default
  jwt_query_param: access_token    # query key for browsers (also accepts Authorization: Bearer)
  sse_enabled: true                 # GET {path_prefix}/sse/:channel
  allow_origins: []                 # empty = permissive CheckOrigin; set explicit origins in production
```

`broadcast.enabled: true` allocates `app.Broadcast` at bootstrap and registers routes. **`jwt_secret`** must be set for **private-** / **presence-** channels (same HS256 rules as `security.JWTMiddleware`).

### Channel naming

| Prefix | Auth | Notes |
|--------|------|--------|
| *(none)* / arbitrary public names | none | Anyone can subscribe. |
| `private-…` | Valid JWT on connect | Optional pattern `private-user-{sub}` restricts the channel to that JWT `sub`. |
| `presence-…` | Valid JWT | Same optional `presence-user-{sub}` pattern; includes **member list** in `pusher_internal:subscription_succeeded` and **`pusher_internal:member_added` / `member_removed`** events. |

### Server-side publish

```go
// app.Broadcast is *broadcast.Hub (set when broadcast.enabled is true).
_ = app.Broadcast.PublishJSON("public-news", "ArticlePublished", map[string]any{"id": 42})
```

`PublishJSON` emits `{ "event", "channel", "data" }` to every WebSocket and SSE subscriber on that channel.

### Presence helper

```go
ids := app.Broadcast.OnlineOn("presence-room")
```

Returns current **JWT `sub`** values tracked for that presence channel (in-process only; not shared across instances).

### Protocol (WebSocket)

After upgrade, the server sends **`pusher:connection_established`** with a **`socket_id`**. Clients subscribe with:

```json
{ "event": "pusher:subscribe", "data": { "channel": "public-news" } }
```

For **presence** channels, include Pusher-style **`channel_data`** JSON: `{"user_id":"…","user_info":{…}}` (optional `user_id` must match JWT `sub` when set).

**Ping:** `{"event":"pusher:ping","data":{}}` → **`pusher:pong`**.

### SSE

`GET /broadcast/sse/my-channel` (or your `path_prefix`) streams **`data:`** lines containing the same JSON envelopes as WebSocket. Use **`?access_token=`** (or your `jwt_query_param`) for **private** / **presence** channels.

---

## Multi-tenancy

ZATRANO supports **tenant resolution** on every HTTP request, optional **row-level** isolation on the generic repository, and **PostgreSQL schema–scoped** migrations for per-tenant DDL.

### Configuration

```yaml
tenant:
  enabled: true
  mode: header              # header | subdomain
  header_name: X-Tenant-ID  # default
  subdomain_suffix: ".app.local"   # required when mode=subdomain: acme.app.local → key acme
  required: false           # true → 400 if tenant missing
  isolation: row            # row | schema (schema sets search_path hint on context)
  row_column: tenant_id     # used by NewTenantAware / TenantScope
  schema_prefix: tenant_    # schema name = prefix + sanitized key
```

When **`tenant.enabled`** is true, **`middleware.ResolveTenant`** runs right after **request-id**. It stores **`tenant.Info`** in:

- **`c.Locals(middleware.LocalsTenant)`** (`*tenant.Info`)
- **`c.Context()`** via **`tenant.WithContext`** so GORM calls with the same `context.Context` see the tenant.

### Row isolation (shared schema)

Use **`repository.NewTenantAware[T](db, "tenant_id")`** instead of **`repository.New[T]`** so every **Find / Paginate / Count / write** adds **`WHERE tenant_id = ?`** from the resolved tenant (numeric **`X-Tenant-ID`** maps to `uint`; non-numeric keys use string equality on the configured column—use a text column such as `tenant_slug` in that case).

Embed **`repository.TenantFK`** on models with a **`tenant_id`** column so **GORM `BeforeCreate`** fills **`TenantID`** from the numeric tenant key when it is zero.

For ad-hoc queries with **`repository.New`**, compose **`repository.TenantScope(ctx, "tenant_id")`** into **`scopes`**.

### Schema isolation (separate PostgreSQL schemas)

1. Resolve tenant with **`isolation: schema`** (sets **`Info.Schema`**, e.g. `tenant_acme`).
2. Obtain a DB handle per request: **`tenant.GormSession(app.DB, c.Context())`** and pass it to **`repository.New[T](scopedDB)`** (no row `tenant_id` filter required if each schema has its own tables).
3. Create schema and migrate:

```bash
zatrano db tenants create-schema --tenant acme
zatrano db tenants migrate --tenant acme
```

These commands build a DSN with **`options=-csearch_path=<schema>,public`** so **golang-migrate** runs **`schema_migrations`** and DDL inside the tenant schema.

**Note:** The in-process hub and cache are not tenant-partitioned automatically; scale-out and RLS policies are application concerns.

---

## Audit / Activity log

ZATRANO provides **model-level activity rows** (who / when / what changed) and optional **HTTP request audit** (method, path, status, latency, user).

### Configuration

```yaml
audit:
  enabled: true
  model_enabled: true       # GORM callbacks → zatrano_activity_logs
  http_enabled: true        # middleware.AuditLog
  http_driver: db           # db | file
  http_file_path: logs/http_audit.jsonl   # required when http_driver is file
```

Run **`zatrano db migrate`** so migration **`000009_audit`** creates **`zatrano_activity_logs`** and **`zatrano_http_audit_logs`**.

### Model activity

1. Register each auditable type once at startup:

```go
import "github.com/zatrano/zatrano/pkg/audit"

audit.RegisterSubject[Product]("products")
```

2. Use the same `context.Context` you pass to GORM, enriched with actor metadata when available:

```go
ctx = audit.WithUser(ctx, userSub)
ctx = audit.WithRequest(ctx, requestID, clientIP)
db.WithContext(ctx).Create(&product)
```

`middleware.AuditLog` already merges **request id**, **IP**, and **JWT `sub`** (when present) into the Fiber request context before your handlers run.

3. **Changes** are stored as an **RFC 6902 JSON Patch** array (shallow object keys; nested values compared by JSON equality). Use **`audit.DiffJSONPatch(oldJSON, newJSON)`** in your own code if needed.

4. **Opt out** for a single call chain: `db.WithContext(audit.Skip(ctx)).Create(...)` skips hooks.

**Soft deletes** implemented as `UPDATE deleted_at` are logged as **updates**, not deletes.

### HTTP audit

**`middleware.AuditLog`** runs after **`AccessLog`** when `audit.enabled` and `audit.http_enabled` are true. It writes either to **`zatrano_http_audit_logs`** (`http_driver: db`) or **append-only JSON lines** (`http_driver: file`). User id resolution order: **`LocalsUserID`** (RBAC middleware) then **JWT claims** `sub`.

### Diff helper

**`audit.DiffJSONPatch`** returns `json.RawMessage` patch bytes suitable for storing in the **`changes`** column on **`zatrano_activity_logs`**.

---

## Full-text search

Define a PostgreSQL **`tsvector`** column in migrations (generated column or trigger-maintained). Query-side helpers live in **`pkg/search`** and return **`repository.Scope`** values:

- **`search.WhereFullTextMatch(regconfig, vectorColumn, userText)`** — no-op when `userText` is empty; otherwise `vectorColumn @@ plainto_tsquery(regconfig, userText)`.
- **`search.OrderByTSRank(regconfig, vectorColumn, userText)`** — descending **`ts_rank_cd`** (no-op when `userText` is empty).

`vectorColumn` must be a **trusted SQL identifier** (constant or allow-list from your code; never pass raw user input as the column name). `regconfig` and the search text are bound as parameters.

Use **`search.postgres_fts_language`** in **`config.Search`** (for example `simple`, `english`, `turkish`) as the default `regconfig` for your app.

### Meilisearch / Typesense

```yaml
search:
  enabled: true
  driver: meilisearch          # or typesense
  default_index_prefix: zatrano_
  meilisearch_url: http://127.0.0.1:7700
  meilisearch_api_key: ""
  # typesense_url / typesense_api_key — required when driver is typesense
  postgres_fts_language: simple
```

When **`search.enabled`** is true, **`core.Bootstrap`** sets **`app.Search`** from **`search.NewClient`**. Physical index / collection names are **`default_index_prefix` + logical name**.

Register bulk importers at startup:

```go
search.RegisterImporter("product", func(ctx context.Context, db *gorm.DB, drv search.Driver) error {
    // load rows, build []search.Document with string IDs, then:
    return drv.UpsertDocuments(ctx, "products", docs)
})
```

Run:

```bash
zatrano search import product
```

Create the target index/collection in the engine first (Meilisearch expects a primary key named **`id`** on documents).

---

## Feature flags

**`pkg/features`** lets you toggle behaviour from **YAML** (`features.definitions`), **PostgreSQL** (`zatrano_feature_flags` via migration **`000010_feature_flags`**), or **both** (`source: both` — a DB row wins over static config for the same key). When **`features.enabled`** is false, every flag resolves to off; **`app.Features`** is still constructed as a no-op.

### Configuration

```yaml
features:
  enabled: true
  source: both              # config | db | both
  definitions:
    - key: beta-ui
      enabled: true
    - key: new-dashboard
      enabled: true
      rollout_percent: 30   # 1..99: requires signed-in user + stable FNV bucket (A/B)
      allowed_roles: [admin, editor]
```

- Non-empty **`allowed_roles`** disables the flag for **anonymous** requests; roles must be present in Fiber Locals via **`middleware.InjectRoles`** (or any code that sets the same keys as `authorize.go`).
- **`rollout_percent`** between 1 and 99 requires a **non-zero numeric user id**; anonymous traffic is treated as out-of-rollout.

### Go API

```go
u := &features.User{ID: 1, Roles: []string{"admin"}}
if app.Features.For(u).IsEnabled(c.Context(), "new-dashboard") {
    // ...
}
// From a request: app.Features.FromFiber(c).IsEnabled(ctx, "beta-ui")
```

### HTTP middleware

```go
import "github.com/zatrano/zatrano/pkg/middleware"

app.Get("/beta", middleware.RequireFeature(app.Features, "beta-ui"), handler)
```

Returns **404** when the flag is off (route behaves as missing). **`server.Mount`** registers **`Features.LocalsMiddleware`** when `features.enabled` is true; use **`features.EvalFromFiber(c)`** to read the same evaluator in handlers.

### View templates

Use **`a.View.ViewData(c, ...)`** so the template root map receives the evaluator binding required by the **`feature`** helper:

```html
{{if feature . "beta-ui"}}
  <p>Beta UI is on</p>
{{end}}
```

The first argument must be the **template root (`.`).** `html/template` cannot call `{{if feature "beta-ui"}}` without that context.

---

## GraphQL

Schema-first GraphQL uses **`gqlgen.yml`** at the module root and **`api/graphql/*.graphqls`**. Generated server code lives under **`pkg/graphql/graph/`**. After editing the schema:

```bash
go run github.com/99designs/gqlgen@v0.17.78 generate
```

The first **`go run github.com/99designs/gqlgen@…`** may take longer while modules download; later runs are cached.

### Configuration

```yaml
graphql:
  enabled: true
  path: /graphql
  playground: true
  playground_path: /playground
```

**`server.Mount`** calls **`graphql.Register`** when **`graphql.enabled`** is true, wiring gqlgen’s `net/http` handler through Fiber’s **`middleware/adaptor`**. The index JSON lists **`graphql`** and, when enabled, **`graphql_playground`**.

### DataLoader (graph-gophers/dataloader)

Each request builds **`graphql.NewLoaders(app)`** and attaches it with **`graphql.WithLoaders`** onto the `context.Context` seen by resolvers. Read it with **`graphql.LoadersFrom(ctx)`** and attach `*dataloader.Loader[...]` fields on **`Loaders`** for batched loads.

### Codegen

```bash
zatrano gen graphql product --module-root .
```

Writes **`api/graphql/product_stub.graphqls`** (example: `extend type Query { product(id: ID!): Product }` plus `type Product { id: ID! }`) and runs **gqlgen generate**. **`--skip-generate`** writes only the `.graphqls` file. Fails if the stub file already exists.

---

## Validation

ZATRANO provides a **generic, struct-tag based validation system** wrapping [`go-playground/validator/v10`](https://pkg.go.dev/github.com/go-playground/validator/v10) with automatic **422 JSON responses** and **i18n-translated error messages**.

### Quick Usage

The primary API is **`zatrano.Validate[T](c)`** — a single generic call that parses the request body, validates struct tags, and returns a structured 422 response on failure:

```go
import "github.com/zatrano/zatrano/pkg/zatrano"

func (h *ProductHandler) Create(c fiber.Ctx) error {
    req, err := zatrano.Validate[CreateProductRequest](c)
    if err != nil {
        return err // 422 JSON response already sent
    }
    // req is valid — use it
    return h.svc.Create(c.Context(), req.Name, req.Email)
}
```

### Form Request Structs

Define your request shapes as plain Go structs with `json` and `validate` tags:

```go
// requests/create_product.go
package requests

type CreateProductRequest struct {
    Name  string `json:"name"  validate:"required,min=2,max=255"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age"   validate:"gte=0,lte=150"`
}
```

```go
// requests/update_product.go
package requests

type UpdateProductRequest struct {
    Name  string `json:"name"  validate:"omitempty,min=2,max=255"`
    Email string `json:"email" validate:"omitempty,email"`
}
```

### Generating Request Structs

Use the CLI to scaffold request stubs automatically:

```bash
# Generate only request structs
zatrano gen request product
# → modules/product/requests/create_product.go
# → modules/product/requests/update_product.go

# gen crud also generates request structs automatically
zatrano gen crud product
# → modules/product/crud_handlers.go      (uses zatrano.Validate[T])
# → modules/product/crud_register.go
# → modules/product/requests/create_product.go
# → modules/product/requests/update_product.go
```

### 422 Error Response Format

When validation fails, the response body follows a consistent JSON structure:

```json
{
  "error": {
    "code": 422,
    "message": "validation failed",
    "details": [
      {
        "field": "Name",
        "tag": "required",
        "message": "This field is required"
      },
      {
        "field": "Email",
        "tag": "email",
        "value": "not-an-email",
        "message": "Must be a valid email address"
      }
    ]
  }
}
```

When i18n is enabled and the request locale is `tr`, messages are automatically translated:

```json
{
  "error": {
    "code": 422,
    "message": "validation failed",
    "details": [
      {
        "field": "Name",
        "tag": "required",
        "message": "Bu alan zorunludur"
      },
      {
        "field": "Email",
        "tag": "email",
        "value": "not-an-email",
        "message": "Geçerli bir e-posta adresi olmalıdır"
      }
    ]
  }
}
```

### i18n Validation Messages

Validation messages are stored in your locale files under the `validation.*` key namespace:

```json
// locales/en.json
{
  "validation": {
    "required": "This field is required",
    "email": "Must be a valid email address",
    "min": "Must be at least {{.Param}} characters",
    "max": "Must be at most {{.Param}} characters"
  }
}
```

The `{{.Param}}` placeholder is replaced with the constraint value (e.g. `min=5` → `"Must be at least 5 characters"`).

**Built-in translated tags:** `required`, `email`, `min`, `max`, `gte`, `lte`, `gt`, `lt`, `len`, `url`, `uri`, `uuid`, `oneof`, `numeric`, `number`, `alpha`, `alphanum`, `boolean`, `contains`, `excludes`, `startswith`, `endswith`, `ip`, `ipv4`, `ipv6`, `datetime`, `json`, `jwt`, `eqfield`, `nefield`.

### Custom Validation Rules

Register custom validation tags with optional i18n support:

```go
import (
    "github.com/go-playground/validator/v10"
    "github.com/zatrano/zatrano/pkg/zatrano"
)

// Register a custom rule
zatrano.RegisterRule("tc_no", func(fl validator.FieldLevel) bool {
    v := fl.Field().String()
    if len(v) != 11 {
        return false
    }
    // ... TC identity number algorithm
    return true
})

// With i18n message key (add "validation.tc_no" to your locale files)
zatrano.RegisterRuleWithMessage("tc_no", tcNoValidator, "validation.tc_no")
```

Then use it in struct tags:

```go
type CitizenRequest struct {
    TCNO string `json:"tc_no" validate:"required,tc_no"`
}
```

### Direct Engine Access

For advanced use cases, access the underlying validator engine:

```go
import "github.com/zatrano/zatrano/pkg/validation"

engine := validation.Default()
engine.Validator() // *validator.Validate from go-playground/validator

// Validate any struct programmatically (without Fiber context)
if verr := engine.ValidateStruct(myStruct, "en"); verr != nil {
    for _, fe := range verr.Errors {
        fmt.Printf("%s: %s\n", fe.Field, fe.Message)
    }
}
```

---

## Authorization (RBAC & Gate/Policy)

ZATRANO provides a **complete authorization system** with two complementary layers: **RBAC** (role-based, DB-backed) for permission checks and **Gate/Policy** (resource-based) for fine-grained instance-level authorization. Both integrate with the **i18n** system for localized 403 error messages.

### RBAC — Role-Based Access Control

Roles and permissions are stored in the database (`roles`, `permissions`, `role_permissions`, `zatrano_user_roles` tables). An in-memory cache avoids DB hits on hot-path permission checks. The `RBACManager` is initialized automatically during bootstrap (when DB is available) and accessible via `app.RBAC`.

```go
import "github.com/zatrano/zatrano/pkg/auth"

// Create roles and permissions
rbac := app.RBAC
rbac.CreateRole(ctx, "admin", "Administrator")
rbac.CreateRole(ctx, "editor", "Content editor")
rbac.CreatePermission(ctx, "posts.create", "Create posts")
rbac.CreatePermission(ctx, "posts.update", "Update posts")
rbac.CreatePermission(ctx, "posts.delete", "Delete posts")

// Assign permissions to roles
rbac.AssignPermissions(ctx, "admin", "posts.create", "posts.update", "posts.delete")
rbac.AssignPermissions(ctx, "editor", "posts.create", "posts.update")

// Assign roles to users
rbac.AssignRoleToUser(ctx, userID, "editor")

// Check permissions
ok, _ := rbac.UserHasPermission(ctx, userID, "posts.create") // true
ok, _ = rbac.UserHasPermission(ctx, userID, "posts.delete")  // false (editor can't delete)
```

**Database migration:** run `zatrano db migrate` — migration `000002_rbac` creates the four required tables with proper indexes and foreign keys.

### Gate / Policy — Resource-Based Authorization

The `Gate` system (accessible via `app.Gate`) allows defining authorization checks for specific actions. Use `Define` for ad-hoc checks or `RegisterPolicy` for structured CRUD policies.

```go
import "github.com/zatrano/zatrano/pkg/auth"

// Ad-hoc gate definition
gate := app.Gate
gate.Define("edit-post", func(c fiber.Ctx, resource any) bool {
    post := resource.(*Post)
    userID, _ := c.Locals(middleware.LocalsUserID).(uint)
    return post.AuthorID == userID
})

// Super-admin bypass (runs before every gate check)
gate.Before(func(c fiber.Ctx, ability string, resource any) *bool {
    roles, _ := c.Locals(middleware.LocalsUserRoles).([]string)
    for _, r := range roles {
        if r == "super-admin" { t := true; return &t }
    }
    return nil // fall through to gate definition
})

// In handlers:
if err := gate.Authorize(c, "edit-post", post); err != nil {
    return err // 403 Forbidden
}
```

**Policy interface:** implement `auth.Policy` for structured CRUD authorization. Generate stubs with `zatrano gen policy`:

```bash
zatrano gen policy post
# → modules/post/policies/post_policy.go
```

The generated policy implements 7 methods: `ViewAny`, `View`, `Create`, `Update`, `Delete`, `ForceDelete`, `Restore`. Register it with the gate:

```go
import "myapp/modules/post/policies"

gate.RegisterPolicy("post", &policies.PostPolicy{})
// Creates: "post.viewAny", "post.view", "post.create", "post.update",
//          "post.delete", "post.forceDelete", "post.restore"
```

### Route-Level Authorization Middleware

The `pkg/middleware` package provides ready-to-use middleware for route-level permission and role checks. All return **403 JSON** with **i18n-aware** error messages.

| Middleware | Description |
|---|---|
| `middleware.Can(rbac, "perm")` | Requires the user to have a specific permission |
| `middleware.CanAny(rbac, "p1", "p2")` | Passes if the user has **any** of the listed permissions |
| `middleware.CanAll(rbac, "p1", "p2")` | Passes only if the user has **all** listed permissions |
| `middleware.HasRole("admin")` | Requires a specific role |
| `middleware.HasAnyRole("admin", "editor")` | Passes if the user has **any** of the listed roles |
| `middleware.GateAllows(gate, "ability")` | Checks a gate ability (without resource) |
| `middleware.InjectRoles(rbac)` | Loads user roles into Locals (place after auth middleware) |

```go
import "github.com/zatrano/zatrano/pkg/middleware"

// After authentication middleware:
app.Use(security.JWTMiddleware(cfg))
app.Use(middleware.InjectRoles(rbac))  // loads roles into context

// Permission-based
app.Get("/admin/users", middleware.Can(rbac, "users.view"), usersHandler)
app.Post("/posts", middleware.Can(rbac, "posts.create"), createPostHandler)
app.Delete("/system", middleware.CanAll(rbac, "system.admin", "system.delete"), handler)

// Role-based
app.Get("/dashboard", middleware.HasAnyRole("admin", "editor"), dashHandler)

// Gate-based
app.Get("/posts", middleware.GateAllows(gate, "post.viewAny"), listPostsHandler)
```

### 403 Error Response Format

When authorization fails, the response follows the standard JSON error shape:

```json
{
  "error": {
    "code": 403,
    "message": "You do not have permission to perform this action.",
    "permission": "posts.delete"
  }
}
```

When i18n is enabled and the request locale is `tr`:

```json
{
  "error": {
    "code": 403,
    "message": "Bu işlemi gerçekleştirme yetkiniz bulunmamaktadır.",
    "permission": "posts.delete"
  }
}
```

### i18n Authorization Messages

Authorization messages are stored under the `auth.*` key namespace in locale files:

```json
{
  "auth": {
    "forbidden": "You do not have permission to perform this action.",
    "unauthorized": "Authentication is required to access this resource.",
    "role_required": "You do not have the required role to access this resource.",
    "permission_required": "You do not have the required permission: {{.Permission}}."
  }
}
```

---

## Cache System

ZATRANO provides a **robust caching layer** with a unified API for **In-Memory** and **Redis** backends. It supports advanced patterns like `Remember`, JSON serialization, tag-based invalidation, and response middleware.

### Drivers

The system automatically chooses the best driver based on your configuration:
- **Redis:** Preferred when `redis_url` is configured. Supports distributed environments and tags.
- **Memory:** Fallback for local development or single-node deployments. Fast, but volatile.

### Basic Usage

Access the cache manager via `app.Cache`:

```go
import "context"

ctx := context.Background()

// Simple storage
app.Cache.Set(ctx, "key", "value", 10 * time.Minute)

// Retrieval
val, ok := app.Cache.Get(ctx, "key")

// Automatic JSON handling
type User struct { Name string }
app.Cache.SetJSON(ctx, "user:1", User{Name: "Alice"}, time.Hour)

var user User
ok, err := app.Cache.GetJSON(ctx, "user:1", &user)
```

### Advanced Patterns

#### `Remember` and `RememberJSON`

The most popular pattern (Laravl-style): returns the cached value if it exists, otherwise computes it via the provided function, caches it, and returns the result.

```go
// Fetch from DB only if not in cache
users, err := app.Cache.RememberJSON(ctx, "users:all", 30*time.Minute, &[]User{}, func() (any, error) {
    return db.FindAllUsers(ctx)
})
```

#### Tags (Redis Only)

Group related keys under tags for bulk invalidation.

```go
// Store under a tag
app.Cache.Tags("users").Set(ctx, "users:1", data, time.Hour)

// Invalidate all keys associated with a tag
app.Cache.Tags("users").Flush(ctx)
```

### Middleware

Cache the entire response of a route at the HTTP level.

```go
import "github.com/zatrano/zatrano/pkg/middleware"

// Cache for 5 minutes
app.Get("/api/v1/stats", middleware.Cache(app.Cache, 5*time.Minute), handler)

// With Tags
app.Get("/api/v1/users", middleware.CacheWithConfig(app.Cache, middleware.CacheConfig{
    TTL:  10 * time.Minute,
    Tags: []string{"users"},
}), handler)
```

### CLI Commands

Clear the cache from the terminal:

```bash
# Clear everything
zatrano cache clear

# Clear specific tags
zatrano cache clear --tag users --tag posts
```

---

## Queue / Job System

ZATRANO provides a **Redis-backed background job queue** with delayed scheduling, automatic retry with exponential backoff, and failed job persistence to PostgreSQL.

### Defining Jobs

Implement the `queue.Job` interface or embed `queue.BaseJob` for sensible defaults:

```go
package jobs

import (
    "context"
    "time"
    "github.com/zatrano/zatrano/pkg/queue"
)

type SendEmailJob struct {
    queue.BaseJob
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

func (j *SendEmailJob) Name() string            { return "send_email" }
func (j *SendEmailJob) Queue() string           { return "emails" }
func (j *SendEmailJob) Retries() int            { return 5 }
func (j *SendEmailJob) Timeout() time.Duration  { return 30 * time.Second }

func (j *SendEmailJob) Handle(ctx context.Context) error {
    // send the email...
    return mailer.Send(ctx, j.To, j.Subject, j.Body)
}
```

Generate job stubs with the CLI:

```bash
zatrano gen job send_email
# → modules/jobs/send_email.go
```

### Dispatching Jobs

```go
// Register job types at startup
app.Queue.Register("send_email", func() queue.Job { return &jobs.SendEmailJob{} })

// Dispatch immediately
app.Queue.Dispatch(ctx, &jobs.SendEmailJob{
    To:      "user@example.com",
    Subject: "Welcome!",
    Body:    "Hello world",
})

// Dispatch with delay (Redis ZADD sorted set)
app.Queue.Later(ctx, 5*time.Minute, &jobs.SendEmailJob{
    To:      "user@example.com",
    Subject: "Follow-up",
})
```

### Worker Process

Start a long-running worker that processes jobs from the queue:

```bash
zatrano queue work
zatrano queue work --queue emails --queue notifications
zatrano queue work --tries 5 --timeout 120s --sleep 5s
```

The worker automatically:
- Polls Redis using BRPOP (FIFO order)
- Migrates delayed jobs (ZADD → LPUSH) every second
- Retries failed jobs with **exponential backoff** (2^attempt seconds)
- Records permanently failed jobs in the `zatrano_failed_jobs` PostgreSQL table
- Recovers from panics inside `Handle()`
- On SIGINT/SIGTERM: stops dequeuing (BRPOP cancelled), waits for the in-flight `Handle()` to finish (single-threaded worker), then exits; the parent `context` is cancelled after `Run` returns so the delayed-job migrator stops

### Failed Jobs

Jobs that exceed their maximum retry count are saved to PostgreSQL with error message, stack trace, and original payload.

```bash
# List failed jobs
zatrano queue failed

# Retry a specific failed job
zatrano queue retry 42

# Retry all failed jobs
zatrano queue retry --all

# Delete all failed job records
zatrano queue flush
```

**Database migration:** run `zatrano db migrate` — migration `000003_failed_jobs` creates the required table.

### Queue Architecture

| Component | Redis Structure | Purpose |
|---|---|---|
| Ready queue | `LIST` (LPUSH/BRPOP) | FIFO job processing |
| Delayed jobs | `SORTED SET` (ZADD) | Time-based scheduling |
| Failed jobs | Database table (`zatrano_failed_jobs` from migrations) | Persistent failure records |

---

## Mail System

ZATRANO provides a **multi-driver mail system** with HTML template support, queue integration for async sending, attachments, and a Mailable pattern for reusable email definitions.

**Notifications:** `core.Bootstrap` registers **`App.Notifications`** (`*notifications.Manager`) with a **`mail`** channel wired to **`App.Mail`**. Built-in **password reset** and **email verification** helpers in **`pkg/auth`** send through **`SendToChannels(..., "mail")`** (plain + optional HTML via `WithData("html", …)` on `notifications.NewNotification`). Add more channels (database, SMS, push) by registering them on the same manager.

### Configuration

```yaml
# config/dev.yaml
mail:
  driver: smtp          # smtp | log (log = dev/testing)
  from_name: "My App"
  from_email: "noreply@myapp.com"
  templates_dir: "views/mails"
  smtp:
    host: smtp.example.com
    port: 587
    username: user
    password: secret
    encryption: tls     # tls | starttls | ""
```

### Sending Emails

```go
import "github.com/zatrano/zatrano/pkg/mail"

// Simple message
app.Mail.Send(ctx, &mail.Message{
    To:      []mail.Address{{Email: "user@example.com", Name: "Alice"}},
    Subject: "Welcome!",
    HTMLBody: "<h1>Hello Alice!</h1>",
})

// With template
app.Mail.SendTemplate(ctx,
    []mail.Address{{Email: "user@example.com"}},
    "Welcome to Our App",
    "welcome",    // views/mails/welcome.html
    "default",    // views/mails/layouts/default.html
    map[string]any{"Name": "Alice"},
)

// Async via queue
app.Mail.Queue(ctx, &mail.Message{
    To:      []mail.Address{{Email: "user@example.com"}},
    Subject: "Newsletter",
    HTMLBody: body,
})
```

### Mailable Pattern

Generate structured, reusable email definitions:

```bash
zatrano gen mail welcome
# → modules/mails/welcome_mail.go
# → views/mails/welcome.html
```

```go
type WelcomeMail struct {
    Name  string
    Email string
}

func (m *WelcomeMail) Build(b *mail.MessageBuilder) error {
    b.To(m.Name, m.Email).
        Subject("Welcome!").
        View("welcome", "default", map[string]any{"Name": m.Name}).
        AttachData("guide.pdf", pdfBytes, "application/pdf")
    return nil
}

// Send synchronously
app.Mail.SendMailable(ctx, &mails.WelcomeMail{Name: "Alice", Email: "alice@example.com"})

// Or async via queue
app.Mail.QueueMailable(ctx, &mails.WelcomeMail{Name: "Alice", Email: "alice@example.com"})
```

### Attachments

```go
msg := &mail.Message{
    To:      []mail.Address{{Email: "user@example.com"}},
    Subject: "Invoice",
    HTMLBody: body,
    Attachments: []mail.Attachment{
        {Filename: "invoice.pdf", Content: pdfBytes},
        {Filename: "logo.png", Content: logoBytes, Inline: true},
    },
}
app.Mail.Send(ctx, msg)
```

### Template Preview

Preview email templates in the browser during development:

```bash
zatrano mail preview              # list templates
zatrano mail preview welcome      # preview welcome template
zatrano mail preview welcome --port 3001
```

For full local mail testing, use **Mailpit** or **MailHog** as the SMTP host.

---

## Event / Listener System

ZATRANO provides a **central event bus** (pub/sub) with support for synchronous and asynchronous listeners, queue-backed delivery via `ShouldQueue`, and generators for rapid development.

### Registering Listeners

```go
// In your service provider / bootstrap (e.g. events/event_service_provider.go):

// Sync listener (inline)
app.Events.ListenFunc("user.created", func(ctx context.Context, e events.Event) error {
    log.Println("user created", e)
    return nil
})

// Struct listener
app.Events.Listen("user.created", &listeners.SendWelcomeMailListener{})

// Multiple listeners for one event
app.Events.Subscribe("order.placed",
    &listeners.SendOrderConfirmationListener{},
    &listeners.UpdateInventoryListener{},
)
```

### Firing Events

```go
import "github.com/zatrano/zatrano/pkg/events"

// Define an event
type UserCreatedEvent struct {
    events.BaseEvent
    UserID uint
    Email  string
}
func (e *UserCreatedEvent) Name() string { return "user.created" }

// Fire synchronously (blocks until all sync listeners complete)
app.Events.Fire(ctx, &UserCreatedEvent{UserID: 1, Email: "alice@example.com"})

// Fire asynchronously (goroutines, errors only logged)
app.Events.FireAsync(ctx, &UserCreatedEvent{UserID: 1, Email: "alice@example.com"})
```

### Async Listeners via Queue (`ShouldQueue`)

Implement `ShouldQueue` to dispatch a listener as a queue job:

```go
type SendWelcomeMailListener struct{}

func (l *SendWelcomeMailListener) Handle(ctx context.Context, event events.Event) error {
    // runs in a background worker
    return nil
}

func (l *SendWelcomeMailListener) Queue() string { return "events" }   // queue name
func (l *SendWelcomeMailListener) Retries() int  { return 3 }
```

When `ShouldQueue` is implemented and a queue is configured (Redis), the listener is automatically dispatched via the Queue system instead of running inline.

### Generator

```bash
zatrano gen event user_created
# → modules/events/user_created_event.go

zatrano gen listener send_welcome_mail
# → modules/listeners/send_welcome_mail_listener.go  (sync)

zatrano gen listener send_welcome_mail --queued
# → modules/listeners/send_welcome_mail_listener.go  (ShouldQueue / async)
```

### Event Service Provider

Centralise all listener registrations in one place:

```go
// modules/events/event_service_provider.go
package myevents

import (
    "github.com/zatrano/zatrano/pkg/core"
    "myapp/modules/listeners"
)

// Register wires all event listeners. Call from main or bootstrap.
func Register(app *core.App) {
    app.Events.Listen("user.created", &listeners.SendWelcomeMailListener{})
    app.Events.Listen("order.placed", &listeners.SendOrderConfirmationListener{})
}
```

---

## Repository / Data

ZATRANO provides a **generic repository pattern** over GORM to standardize data access, enforce reusable query scopes, and automate common tasks like pagination and soft-deleting.

### Base Model & Generic Repository

Use the `repository.Model` to get ID, timestamps, and standard soft-delete behavior out of the box.

```go
import "github.com/zatrano/zatrano/pkg/repository"

type User struct {
    repository.Model
    Name  string
    Email string
}

// In your service layer:
repo := repository.New[User](app.DB)

// Create
repo.Create(ctx, &User{Name: "Alice", Email: "alice@example.com"})

// Soft Delete
repo.DeleteByID(ctx, 1)

// Restore
repo.Restore(ctx, 1)
```

### Chainable Scopes

Build complex queries without leaking GORM internals into your handlers. 

```go
// Pre-defined scopes
scopes := repository.Scopes(
    repository.Active(),
    repository.Where("email LIKE ?", "%@example.com"),
    repository.PreloadAll(), // Eager load to prevent N+1
    repository.OrderBy("created_at DESC"),
    repository.Limit(10),
)

users, _ := repo.FindAll(ctx, scopes...)
```

### Pagination

Pagination is built-in and standardizes responses. `repo.Paginate` returns a `Page[T]` containing items and normalized metadata.

```go
opts := repository.PaginateOpts{Page: 1, PerPage: 15}

page, _ := repo.Paginate(ctx, opts, repository.Active())

// page.Items (your data)
// page.Pagination.Total, page.Pagination.CurrentPage, etc.

// Get HTML pagination links for rendering in templates
links := page.Pagination.Links("/users", "&sort=desc")
```

---

### Internationalization (`i18n`)

Application UI copy lives in **JSON** files under **`locales_dir`**, one file per locale: **`{locales_dir}/{tag}.json`** (e.g. `locales/en.json`). Nested objects are flattened to **dot keys** (`app.welcome`).

- **Config:** `i18n.enabled`, `i18n.default_locale`, `i18n.supported_locales`, `i18n.locales_dir`, optional `i18n.cookie_name` (default `zatrano_lang`), `i18n.query_key` (default `lang`). When **`i18n.enabled`** is true, **`locales_dir`** must exist on disk (validated at config load).
- **Resolution order:** query (`?lang=`), cookie, **`Accept-Language`**, then **`default_locale`**.
- **Handlers:** `import "github.com/zatrano/zatrano/pkg/i18n"` — **`i18n.T(c, "app.welcome")`** for static strings; **`i18n.Tf(c, "app.hello_user", map[string]any{"Name": userName})`** (or any struct) for **`text/template`** placeholders such as **`{{.Name}}`** in JSON. For **`map`** data, simple `{{.Field}}` segments are rewritten automatically; use **`Bundle.Format(locale, key, data)`** without Fiber. If i18n is off, **`T`** / **`Tf`** return the key unchanged ( **`Tf`** also returns **`nil`** error).
- **GET /** includes an **`i18n`** object (`enabled`, and when on: `default_locale`, `supported_locales`, **`active_locale`** for the current request).
- **Validation messages** are automatically resolved from `validation.*` keys when i18n is enabled (see [Validation](#validation)).

---

## View / Template System

ZATRANO ships a first-class server-rendered template engine built on top of Go's `html/template`, with layout inheritance, reusable component partials, a rich form-builder helper set, session-backed flash messages, old-input repopulation, and a Vite/esbuild-aware asset pipeline.

### Layout Inheritance (`{{extends}}` / `{{block}}`)

Child views declare their parent layout on the **first non-blank line**:

```html
{{extends "layouts/app"}}

{{block "title"}}Dashboard{{end}}

{{block "content"}}
<h1>Welcome, {{.User.Name}}</h1>
{{end}}
```

The engine reads the `{{block "name"}}default{{end}}` declarations in the layout and replaces them with the child's overrides. Any block the child omits renders its default content from the layout.

**Built-in layout blocks in `layouts/app`:**

| Block | Purpose |
|-------|---------|
| `title` | `<title>` content |
| `head` | Extra `<head>` elements (meta, styles) |
| `body_class` | CSS classes on `<body>` |
| `header` | Top navigation bar |
| `nav` | Navigation links inside the header |
| `content` | Main page content (**required**) |
| `footer` | Page footer |
| `scripts` | Extra `<script>` tags before `</body>` |

### Component System (`views/components/`)

Components are plain `.html` files under `views/components/`. They are auto-discovered and registered as named templates at boot — no manual imports needed.

```html
{{/* Inline alert */}}
{{template "components/alert" (dict "Type" "success" "Message" "Saved!")}}

{{/* Form input with validation error and old value */}}
{{template "components/form-input" (dict
  "Type"     "email"
  "Name"     "email"
  "Label"    "Email Address"
  "Value"    (old "email" .Old)
  "Required" true
  "Error"    (index .Errors "email")
)}}
```

**Built-in components:**

| Component | Description |
|-----------|-------------|
| `components/alert` | Coloured alert box (`success`, `error`, `warning`, `info`) with SVG icon |
| `components/button` | `<button>` with variant (`primary`, `secondary`, `danger`, `ghost`) |
| `components/form-input` | `<input>` wrapped in label + error/hint |
| `components/form-select` | `<select>` with option list, label, error |
| `components/form-textarea` | `<textarea>` with label, rows, error |
| `components/csrf` | Hidden `<input name="_csrf">` |
| `components/pagination` | Offset-based pagination links |
| `partials/flash-messages` | Renders all queued flash messages |

### Form Builder

Template functions for building HTML forms without raw HTML:

```html
{{form_open "/users" "POST"}}
  {{csrf_field .CSRF}}

  {{input "text" "name" (old "name" .Old) `class="form-control"`}}
  {{textarea "bio" (old "bio" .Old) `rows="4"`}}
  {{select "role" .FormRole (slice (arr "admin" "Admin") (arr "user" "User"))}}
  {{checkbox "active" "1" true}}

  <button type="submit">Save</button>
{{form_close}}
```

| Helper | Signature | Output |
|--------|-----------|--------|
| `form_open` | `action method [attrs...]` | `<form ...>` |
| `form_close` | — | `</form>` |
| `csrf_field` | `token` | `<input type="hidden" name="_csrf" value="...">` |
| `input` | `type name value [attrs...]` | `<input ...>` |
| `textarea` | `name value [attrs...]` | `<textarea>...</textarea>` |
| `select` | `name selected [][2]string [attrs...]` | `<select>...</select>` |
| `checkbox` | `name value checked [attrs...]` | `<input type="checkbox" ...>` |

### Flash Messages

Set a flash message before redirecting; it is available in the next request and then cleared automatically.

```go
import "github.com/zatrano/zatrano/pkg/view/flash"

// In a handler:
flash.Set(c, flash.Success, "Record saved successfully.")
flash.Set(c, flash.Error,   "Something went wrong.")
return c.Redirect("/dashboard")
```

Flash types: **`flash.Success`**, **`flash.Error`**, **`flash.Warning`**, **`flash.Info`**.

In your layout or view, `{{template "partials/flash-messages" .}}` renders all queued messages. Flash data is injected automatically by `a.View.ViewData(c)`.

### Old Input Helper (`{{old}}`)

After a failed form validation, persist the user's input so the form repopulates:

```go
// In the handler — before redirecting back:
flash.SetOld(c, map[string]string{
    "email": c.FormValue("email"),
    "name":  c.FormValue("name"),
})
flash.Set(c, flash.Error, "Validation failed.")
return c.Redirect("/users/create")
```

In the template:

```html
{{input "email" "email" (old "email" .Old) `placeholder="you@example.com"`}}
```

Or via the callable function injected by `ViewData`:

```html
<input type="email" name="email" value="{{call .OldFn "email"}}">
```

### Asset Helper (`{{asset}}`)

`{{asset "path"}}` resolves a versioned URL for a static file. Resolution order:

1. **Vite/esbuild manifest** — returns the hashed filename (e.g. `app-a1b2c3.js`)
2. **MD5 file hash** — appends `?v=<hash>` for files in `PublicDir`
3. **Plain URL** — falls back to `PublicURL/path`

```html
{{/* Returns e.g. /public/css/app.css?v=1a2b3c4d */}}
<link rel="stylesheet" href="{{asset "css/app.css"}}">

{{/* Returns full <link> or <script> tag */}}
{{assetLink "css/app.css"}}
{{assetScript "js/app.js"}}
```

### Vite / esbuild Integration

Point `view.asset.vite_manifest` at your build output manifest. After `vite build`, `{{asset}}` resolves hashed filenames automatically. During development, enable HMR:

```yaml
# config/dev.yaml
view:
  dev_mode: true
  asset:
    vite_dev_url: http://localhost:5173
    vite_manifest: public/build/.vite/manifest.json
```

```html
{{/* Injects Vite client + module entry point in dev, hashed tags in prod */}}
{{viteHead "src/main.ts"}}
```

In **dev mode** with `ViteDevURL` set, `{{viteHead "src/main.ts"}}` outputs:

```html
<script type="module" src="http://localhost:5173/@vite/client"></script>
<script type="module" src="http://localhost:5173/src/main.ts"></script>
```

After `vite build`, the same call reads `manifest.json` and emits hashed `<link>` + `<script>` tags.

### Using the View Engine in Handlers

```go
// In any handler that has access to *core.App:
func (h *Handler) Show(c fiber.Ctx) error {
    data := h.app.View.ViewData(c, fiber.Map{
        "Title": "Dashboard",
        "User":  user,
    })
    // ViewData automatically injects: Flash, Old, OldFn, CSRF, Path, Method
    return c.Render("dashboard/show", data)
}
```

Register the renderer as Fiber's view engine in your `RegisterRoutes` hook (generated apps do this automatically):

```go
fiberApp := core.NewFiber(a)
fiberApp.Set("Views", a.View)   // implements fiber.Views
```

### Template Helper Reference

| Function | Example | Description |
|----------|---------|-------------|
| `asset` | `{{asset "app.css"}}` | Versioned asset URL |
| `assetLink` | `{{assetLink "app.css"}}` | `<link rel="stylesheet">` tag |
| `assetScript` | `{{assetScript "app.js"}}` | `<script defer>` tag |
| `viteHead` | `{{viteHead "src/main.ts"}}` | Vite entry point tags (dev + prod) |
| `old` | `{{old "email" .Old}}` | Previous form value |
| `csrf_field` | `{{csrf_field .CSRF}}` | Hidden CSRF input |
| `dict` | `{{dict "K" "V"}}` | Build `map[string]any` inline |
| `safe` | `{{safe .HTML}}` | Mark string as `template.HTML` |
| `safeURL` | `{{safeURL .URL}}` | Mark string as `template.URL` |
| `nl2br` | `{{nl2br .Text}}` | Replace newlines with `<br>` |
| `default` | `{{default "—" .Val}}` | Fallback when value is zero |
| `upper` / `lower` / `title` | `{{upper .Name}}` | String case transformers |
| `json` | `{{json .Data}}` | JSON-encode a value |
| `hasKey` | `{{hasKey .Map "key"}}` | Check map contains key |
| `concat` | `{{concat "a" "b"}}` | Concatenate strings |
| `iterate` | `{{range iterate 5}}` | Iterate 0…n-1 |

### Code Generation (`zatrano gen view`)

```bash
# Scaffold index.html + show.html for a module:
zatrano gen view post

# Also generate create.html + edit.html with form scaffolding:
zatrano gen view post --with-form

# Use a custom layout:
zatrano gen view post --with-form --layout layouts/admin

# Dry run — print files without writing:
zatrano gen view post --with-form --dry-run
```

### Configuration (`view.*`)

```yaml
view:
  root: views              # template root directory (default: views)
  extension: .html         # template file extension (default: .html)
  components_dir: components  # sub-dir for component partials
  layouts_dir: layouts     # sub-dir for layout templates
  dev_mode: true           # disable caching; hot-reload on every request

  asset:
    public_dir: public          # filesystem path for static assets
    public_url: /public         # URL prefix for asset URLs
    vite_manifest: ""           # path to Vite/esbuild manifest.json
    vite_dev_url: ""            # Vite dev server URL (HMR in dev_mode)
```

---

## Scheduled tasks

**`pkg/schedule`** registers named cron jobs (standard 5-field spec via **`robfig/cron/v3`**). Use the fluent builder from **`schedule.Call(fn)`** (`EveryMinute`, `Hourly`, `Daily`, `Weekly`, `Monthly`, **`At("HH:MM")`** on day-based presets, **`WithSpec`**, **`Name`**, **`Description`**). Call **`Register()`** at process startup (typically `init()` in `modules/.../schedule.go`).

**Overlap protection:** **`WithoutOverlap()`** serializes a task across replicas using a Redis lock (**`SkipOverlap`** + **`redis_url`** required for **`zatrano schedule run`**).

```bash
zatrano schedule list
zatrano schedule run --env dev --config-dir config
```

The **`schedule run`** command loads config, opens Redis when any task uses overlap protection, starts the internal cron scheduler, and blocks until **SIGINT** / **SIGTERM**.

---

## Admin dashboard

Two related surfaces:

1. **Framework admin (`pkg/admin`)** — When **`admin.enabled: true`**, **`server.Mount`** calls **`admin.Register`**, which mounts **HTML** pages under **`admin.path_prefix`** (default **`/admin`**): dashboard, **metrics** (includes configured Redis queue depths), **logs** (tail of **`admin.log_file`**), and **log download**. If **`admin.secret`** is set, every request must send **`X-Admin-Key`** or **`?admin_key=`**. In **`env: prod`**, **`admin.secret`** is **required** when admin is enabled.

```yaml
admin:
  enabled: true
  path_prefix: /admin
  secret: "change-me"       # required in prod
  log_file: storage/logs/app.log
  queue_names: [default, mails, events]
```

2. **Scaffolded module admin (`zatrano gen admin <module>`)** — Generates **`admin_handlers.go`**, **`admin_register.go`**, and **`views/admin/<module>/index.html`** for a list UI (requires an existing **`zatrano gen module`** tree and model). Use **`RegisterAdmin`** from your wire registration so routes mount under the same prefix helpers from **`pkg/admin`**.

---

## REST API helpers

**`pkg/api`** standardizes JSON APIs beyond raw Fiber handlers.

- **`api.Success` / `api.Created` / `api.JSON`** — `{ "data", "meta?", "links?" }` envelope.
- **`api.CursorPage`**, **`api.CursorPaginateOpts`**, **`EncodeCursor` / `DecodeCursor`** — keyset pagination metadata.
- **`api.NewVersionManager(cfg)`** — **`Group`**, **`Middleware`** (**`X-API-Version`**), driven by **`api.*`** in config.
- **`api.NewKeyManager(db)`**, **`api.GenerateKey`**, **`api.ValidateKey`** — **`api_keys`** table helpers for CLI **`zatrano api-key …`** and application rotation flows.

---

## Outbound HTTP client

**`pkg/http`** provides **`http.NewClient()`** for calling external JSON services: **`WithToken`**, **`WithHeader`**, **`WithTimeout`**, **`WithRetry`**, and method builders **`Get` / `Post` / `Put`** returning a fluent **`Request`** with JSON decode helpers. **`http.Fake(handler)`** swaps the transport for deterministic tests.

```go
c := http.NewClient().WithToken(os.Getenv("REMOTE_TOKEN")).WithTimeout(10 * time.Second)
var out RemoteDTO
_ = c.Get("https://api.example.com/v1/me").Into(&out)
```

---

## Notifications

**`pkg/notifications`** defines **`Notification`**, **`Channel`**, and **`Manager`** (`Register`, `Send`, `SendToChannels`). **`core.Bootstrap`** wires a default **`mail`** channel to **`App.Mail`**; add database, SMS, or push channels by **`Register`**-ing additional implementations (see **`pkg/notifications`** drivers). Use **`notifications.NewNotification`** / **`BaseNotification`** for simple payloads.

```go
n := notifications.NewNotification("Welcome", "Thanks for signing up.", user.Email)
_ = app.Notifications.SendToChannels(ctx, n, "mail")
```

**`zatrano gen notification <Name>`** scaffolds a typed notification stub. Password reset / e-mail verification in **`pkg/auth`** already send through **`SendToChannels(..., "mail")`**.

---

## HTTP testing helpers

**`pkg/testing`** supports black-box HTTP tests and database isolation.

- **`testing.NewHTTPClient(fiberApp)`** — **`Get` / `Post` / `Put` / `Patch` / `Delete`** with JSON bodies, optional bearer token, **`HTTPResponse`** helpers (**`Status`**, **`JSON`**, assertions).
- **`testing.NewTestSuite(db)`** — **`SetupTest` / `TeardownTest`** pair wraps each test in a rolled-back GORM transaction via **`GetDB()`**.

---

## Storage / File Management

**`pkg/storage`** manages named **disks** (default **`local`**, optional **`private`**, or **S3-compatible** **`s3` / `minio` / `r2`** via **`storage.Register`** / config map). Features include **signed / temporary URLs**, **image helpers**, and **`storage.ServeFileMiddleware`** for guarded downloads.

Typical wiring reads **`storage.*`** from **`config.Config`** (see `config/examples/dev.yaml`), builds a **`storage.Manager`**, obtains a **`Driver`** with **`manager.Disk("local")`** (or your disk name), then uses **`Put`**, **`Get`**, **`TemporaryURL`**, etc.; expose downloads with **`storage.ServeFileMiddleware`** when files are private.

```bash
zatrano storage:link
zatrano storage:clear local --force
```

Use **`storage:link`** to symlink **`storage/app/public`** into **`public/storage`** for browser-visible uploads.

---

## Configuration

- **`.env`**, **`config/{env}.yaml`**, **environment variables** (nested keys use underscores, e.g. `SECURITY_JWT_SECRET`). For **lists** (e.g. multiple CORS origins or **`supported_locales`**), prefer **YAML**; env overrides for slices vary by shell.
- Key fields: `migrations_source`, `migrations_dir`, `seeds_dir`, `openapi_path`, **`http.*`**, **`i18n.*`**, **`admin.*`**, **`storage.*`**, `security.*`, `oauth.*`, **`api.*`**, **`broadcast.*`**, **`graphql.*`**, … (see `config/examples/dev.yaml`).

### Database migrations (SQL)

Supported engines, **`database_driver`**, and **`database_url`** are described in **[Database](#database)**.

- **`migrations_source`:** **`embed`** (default) — versioned `*.up.sql` / `*.down.sql` live under **`pkg/migrations/sql/<driver>/`** (`postgres`, `mysql`, `sqlite`, `sqlserver`). `zatrano db migrate` uses **golang-migrate** with an **`embed`/`iofs`** source and the same driver you set with **`database_driver`**.
- **`file`** — read migrations from **`migrations_dir`** on disk (typical for **`zatrano new`** / scaffolded apps, which set `migrations_source: file` and ship starter SQL under `migrations/`).
- **`--migrations <dir>`** on **`db migrate`**, **`db rollback`**, or **`db tenants …`** forces **file** mode from that directory (ignores embed for that invocation).
- **`zatrano gen model`** writes new **`.up.sql` / `.down.sql`** stubs under **`pkg/migrations/sql/postgres/`** only; copy or adapt for other drivers if you rely on **embed** for those engines.
- The repo-root **`migrations/`** folder is optional disk staging; see **`migrations/README.md`** when using **`file`** mode.
- Debug: **`zatrano config print`** (full dump, redacted) or **`zatrano config print --paths-only`** (env, cwd, profile path, dirs — safe to paste in chat).
- CI: **`zatrano config validate -q`** (fast YAML/env checks), then **`zatrano openapi validate --merged`**, or **`zatrano verify`** for the full gate (see Development).

---

## Development

```bash
go test ./... -count=1
go fmt ./...
go vet ./...
golangci-lint run   # when installed
```

**One-shot gate:** `zatrano verify` (or **`make verify`** on POSIX) runs `vet`, `test`, and merged OpenAPI validation. **`make verify-race`** / **`zatrano verify --race`** before release builds (slower; catches data races). **`make config-validate`** mirrors **`zatrano config validate`**.

**Live reload:** install [Air](https://github.com/air-verse/air), then `air` (uses `.air.toml`). On Windows the binary is `./tmp/main.exe`.

**Merged OpenAPI file:** `make openapi-export` (POSIX Make) or `go run ./cmd/zatrano openapi export --output api/openapi.merged.yaml`.

**Environment check:** `zatrano doctor` prints config summary, **`http` middleware** (CORS, rate limit, timeout, body limit) and a pointer to **`config print --paths-only`**, **OAuth** when enabled, **`pg_dump` / `pg_restore` / `psql`** PATH resolution for backup/restore, then connectivity probes.

**Generate code:** `gen module` / `gen crud` patch **`zatrano:wire:*`** markers and run **`go fmt`** on the wire file. **`gen wire`** does the same patch without regenerating `modules/` (e.g. after **`--skip-wire`**). **`gen request`** generates form request struct stubs independently. **`gen policy`** generates an `auth.Policy` implementation with CRUD methods (ViewAny, View, Create, Update, Delete, ForceDelete, Restore). **Apps:** `internal/routes/register.go`. **Framework checkout:** `pkg/server/register_modules.go`.

**Embedding the server:** `server.Mount(app, fiberApp, server.MountOptions{RegisterRoutes: …})`; `zatrano.StartOptions.RegisterRoutes` passes through for generated apps.

---

## Documentation

- **English:** this file (`README.md`)
- **Türkçe:** [`README.tr.md`](README.tr.md)

Keep both in sync when adding or changing features.

---

## Contributing

Issues and PRs are welcome. For any behavior or CLI change, update **both** `README.md` and `README.tr.md` in the same change.

---

## License

To be determined.
