# Changelog

All notable changes to ZATRANO are documented in this file.

## 0.1.1 - 2026-08-06

### Added

- Nested markdown discovery for the documentation engine (`core/docs`)
- Sidebar navigation via optional `navigation.json`
- Previous/next page neighbors and custom `ViewRenderer` support for docs routes

## 0.1.0 - 2026-08-05

Initial release.

### Added

- Go module `github.com/zatrano/framework`
- Application skeleton: `app/`, `bootstrap/`, `config/`, `routes/`, `views/`, `database/`, `storage/`, `cmd/zatrano`
- HTTP request/response layer and router
- Controller registration helpers
- View engine (layouts, components, directives)
- Config loading via `.env` and `config/`
- Middleware pipeline (CSRF, CORS, security headers, throttle, and related helpers)
- Database migrations, schema builder, and ORM (`core/orm`)
- Logging, session, cache, queue, mail, validation, and auth packages in `core/`
- CLI entrypoint (`serve`, `migrate`, `make:*`, `about`, …)
- Welcome web route and home controller
- GitHub Actions CI (tests, static analysis, coding style)
- MIT license (Serhan KARAKOÇ)
