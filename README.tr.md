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
[![GORM](https://img.shields.io/badge/GORM-632CA6?style=for-the-badge)](https://gorm.io/)
[![Zap](https://img.shields.io/badge/Zap-yapılandırılmış%20log-121212?style=for-the-badge)](https://github.com/uber-go/zap)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-6BA539?style=for-the-badge&logo=openapiinitiative&logoColor=white)](https://www.openapis.org/)
[![GraphQL](https://img.shields.io/badge/GraphQL-E10098?style=for-the-badge&logo=graphql&logoColor=white)](https://graphql.org/)
[![gqlgen](https://img.shields.io/badge/gqlgen-311C87?style=for-the-badge)](https://github.com/99designs/gqlgen)
[![golang-migrate](https://img.shields.io/badge/golang--migrate-SQL-00599C?style=for-the-badge)](https://github.com/golang-migrate/migrate)
[![Cobra CLI](https://img.shields.io/badge/Cobra-CLI-7E43B6?style=for-the-badge)](https://github.com/spf13/cobra)
[![Viper](https://img.shields.io/badge/Viper-yapılandırma-273F5B?style=for-the-badge)](https://github.com/spf13/viper)
[![AWS SDK](https://img.shields.io/badge/AWS-S3%20SDK-232F3E?style=for-the-badge&logo=amazonaws&logoColor=white)](https://aws.amazon.com/sdk-for-go/)
[![OAuth2](https://img.shields.io/badge/OAuth2-x--oauth2-4285F4?style=for-the-badge)](https://pkg.go.dev/golang.org/x/oauth2)

**Çoklu veritabanı (GORM + `zatrano db migrate`):** PostgreSQL (varsayılan) · MySQL · SQLite · SQL Server — `database_driver` + `database_url`; sürücü başına gömülü SQL: [`pkg/migrations/sql/`](pkg/migrations/sql/).

</div>

---

**ZATRANO**, Go için **bir web framework’ü değil**; uçtan uca bir **backend platform ekosistemidir**: HTTP çalışma zamanı, güvenlik ve kimlik doğrulama, veri erişimi, önbellek ve kuyruk, e-posta ve olaylar, üretken **kod jeneratörleri** ve **CLI** tek çatı altında birleşir. Amaç; **modüler monolit** veya servis tarzı backend’leri **aynı dilde** üretmek, yapılandırmak, migrate etmek ve işletmektir.

Bu depo bilinçli olarak **“Fiber + birkaç middleware”** seviyesinde tutulmaz: üretim öncesi kütüphaneler (Fiber, GORM, Redis, Zap, OpenAPI, gqlgen, golang-migrate, …) üzerine oturan **ürün katmanı** ve **araç zinciri** sunar.

- **Modül yolu:** `github.com/zatrano/zatrano`
- **Go:** 1.25+
- **Çekirdek yığın:** Fiber v3, **PostgreSQL / MySQL / SQLite / SQL Server** (`database_driver` + GORM), Redis, GORM, Zap, **golang-migrate** (sürücüye özel gömülü SQL: `pkg/migrations/sql/<sürücü>/`, `migrations_source`), OpenAPI; isteğe bağlı GraphQL (gqlgen), AWS S3 SDK, OAuth2 (`x/oauth2`)

> **Durum:** aktif geliştirme. Genel Go API’leri **`pkg/`** altındadır; platform üzerine kurulan uygulamalar bu sözleşmeleri import eder.

### Geliştirici

**Serhan KARAKOÇ** — [github.com/serhankarakoc](https://github.com/serhankarakoc)

---

## İçindekiler

- [Özellikler](#özellikler-yol-haritası)
- [Dizilim](#dizilim-pkg-ve-internal)
- [Gereksinimler](#gereksinimler)
- [Kurulum](#kurulum)
- [Hızlı Başlangıç](#hızlı-başlangıç)
- [CLI Komutları](#cli-komutları)
- [HTTP Rotaları & Uluslararasılaştırma (i18n)](#http-şimdilik)
- [Validation (Doğrulama)](#validation-doğrulama)
- [Yetkilendirme (RBAC & Gate/Policy)](#yetkilendirme-rbac--gatepolicy)
- [Cache Sistemi (Önbellek)](#cache-sistemi-önbellek)
- [Kuyruk / Job Sistemi](#kuyruk--job-sistemi)
- [Mail Sistemi (E-posta)](#mail-sistemi-e-posta)
- [Event / Listener Sistemi](#event--listener-sistemi)
- [Broadcasting / WebSocket](#broadcasting--websocket)
- [Multi-tenancy](#multi-tenancy)
- [Audit / Activity log](#audit--activity-log)
- [Full-text search](#full-text-search)
- [Feature flags](#feature-flags)
- [GraphQL](#graphql)
- [Repository / Veri Sistemi](#repository--veri-sistemi)
- [Storage / Dosya Yönetimi](#storage--dosya-yönetimi)
- [View / Template Sistemi](#view--template-sistemi)
- [Yapılandırma](#yapılandırma)
  - [Veritabanı migrasyonları (SQL)](#veritabanı-migrasyonları-sql)
- [Geliştirme](#geliştirme)
- [Dokümantasyon](#dokümantasyon)
- [Katkı](#katkı)
- [Lisans](#lisans)

---

## Özellikler (yol haritası)

| Alan | Plan |
|------|------|
| Mimari | Modüler çekirdek + takılabilir modüller (modüler monolith) |
| Katmanlar | Handler → Service → Repository (zorunlu; hepsinde base) |
| Web | Fiber HTML şablonları, CSRF, **validation** (`go-playground/validator`), flash, **CORS**, **rate limit**, **i18n** (JSON çeviriler), **cache** (Memory/Redis), güvenlik başlıkları, gzip, static |
| **View Motoru** | Layout kalıtımı (`{{extends}}`), blok/bölüm sistemi, bileşen partial'ları, form builder yardımcıları, flash mesajlar, eski girdi yeniden doldurma, versiyonlanmış asset URL'leri, Vite/esbuild manifest entegrasyonu, HMR dev sunucu proxy |
| API | REST + **OpenAPI 3** (`api/openapi.yaml`, `/docs`, `/openapi.yaml`), **Resource/Transformer** (model→JSON, hassas alanları gizle, ilişkileri şekillendir), **Standart response zarfı** ({data, meta, links}, JSON:API uyumlu), **Cursor pagination** (büyük dataset'ler için keyset), **Throttle** (kullanıcı/JWT subject rate limiting, Redis sayaçları), **API key yönetimi** (api_keys tablosu, middleware, rotasyon), **Versioning yöneticisi** (v1/, v2/ otomatik gruplar, config'den prefix) |
| Kimlik | **Oturum (Redis) + CSRF**; `/api/v1/private/*` için **JWT**; **OAuth2** (Google/GitHub) tarayıcı girişi; **RBAC** (rol→izin, DB destekli); **Gate/Policy** (kaynak bazlı yetkilendirme); **Şifre sıfırlama** / **e-posta doğrulama** (işlemsel e-posta **`pkg/notifications`** → **`mail`** kanalı, `App.Notifications` + `App.Mail`); **Brute Force Koruması** (IP+username rate limiting, Redis); **TOTP 2FA** (Google Authenticator uyumlu, QR kod üretimi); **Oturum Yönetimi** (aktif oturumları listele/sonlandır, cihaz bilgisi); **JWT Refresh Token'ları** (token rotasyonu, refresh token tablosu) |
| **Test Altyapısı** | **HTTP test client** (Fiber.Test() sarmalama, Get/Post/WithToken, AssertStatus/AssertJSON), **Database factory** (gofakeit tabanlı test verisi üretimi, gen factory), **Transaction rollback** (TestSuite struct, SetupTest/TeardownTest), **In-memory cache driver** (Redis gerektirme), **Mail fake** (mailleri bellekte tut, gönderildiğini assert et), **Queue fake** (dispatch edilen job'ları assert et) |
| Veri | **Generic Repository** deseni, otomatik soft-delete, **zincirleme Scope'lar**, Offset tabanlı sayfalama |
| VT / Ops | **PostgreSQL · MySQL · SQLite · SQL Server** (`database_driver` + GORM); **`db migrate` / `rollback`** (varsayılan **embed** SQL **`pkg/migrations/sql/<sürücü>/`**), **`seed`**, **`db backup` / `restore`** (Postgres istemci araçları) |
| **Depolama** | **Yerel / S3 / MinIO / Cloudflare R2** sürücüleri, **imzalı URL'ler**, **resim işleme** (yeniden boyutlandırma, kırpma, küçük resim), **Fiber middleware**, genel + özel diskler |
| **HTTP Client** | Zincirleme API ile JSON odaklı HTTP istemcisi; **WithToken**, **WithHeader**, **WithTimeout**, `Get`/`Post`/`Put`, otomatik JSON marshal/unmarshal, 5xx hatalarında retry ve testler için fake transport |
| Kuyruk | **Redis tabanlı** job kuyruğu, geciktirilmiş joblar (ZADD), otomatik retry + üssel geri çekilme, başarısız joblar (veritabanı tablosu `zatrano_failed_jobs`, migration ile) |
| **Zamanlanmış Görevler** | `robfig/cron/v3` sarmalaması, `schedule.Call(fn).Daily().At("08:00")`, `EveryMinute`/`Hourly`/`Daily`/`Weekly`/`Monthly`, Redis overlap kilidi |
| Mail | **SMTP / Log** sürücüleri, HTML şablon + layout desteği, kuyruk entegrasyonu, ek dosya, Mailable deseni |
| Events | **Senkron ve asenkron** event bus, `ShouldQueue` ile kuyruk tabanlı listener, `gen event` + `gen listener` |
| Notifications | Çok kanallı gönderim (Veritabanı, Mail, SMS, Push), okundu/okunmadı takibi, **Twilio / Netgsm sürücüleri**, **FCM / APNs**, `gen notification` |
| **Broadcasting** | **WebSocket hub** (`pkg/broadcast`, kanal yayını, `github.com/fasthttp/websocket` + Fiber v3), **private / presence** kanalları (JWT `sub`), **çevrimiçi listesi** (`Hub.OnlineOn`), **SSE** tek yönlü push, **Pusher uyumlu** kablo formatı (Echo / pusher-js) |
| **Multi-tenancy** | **`ResolveTenant`** ara katmanı (başlık `X-Tenant-ID` veya alt alan adı), **`tenant.FromContext` / Locals**, **satır izolasyonu** (`repository.NewTenantAware` + `TenantScope`), isteğe **`TenantFK`** gömülü yapı; **şema izolasyonu** (`tenant.GormSession` + **`zatrano db tenants`** migrate/rollback/create-schema, PostgreSQL `search_path`) |
| **Audit** | **Model activity** (`zatrano_activity_logs`, GORM callback'leri, `audit.RegisterSubject`, **JSON Patch** farkı `audit.DiffJSONPatch`), **HTTP audit** (`middleware.AuditLog`, `zatrano_http_audit_logs` veya **dosyaya** JSONL), **`audit.WithUser` / `WithRequest`** ile `context.Context` |
| **Full-text search** | **PostgreSQL** `tsvector` / `plainto_tsquery` için **`repository.Scope`** (`pkg/search`), **Meilisearch / Typesense** hafif HTTP sürücüsü, **`zatrano search import <Model>`** ile toplu indeks (`search.RegisterImporter`) |
| **Feature flags** | **`pkg/features`** — YAML ve/veya **`zatrano_feature_flags`** tablosu, kullanıcı + rol + **yüzde rollout** (A/B), **`app.Features.For(user).IsEnabled`**, **`middleware.RequireFeature`**, şablon **`{{if feature . "anahtar"}}`** (`ViewData` ile) |
| **GraphQL** | **gqlgen** şema öncelikli (`api/graphql/*.graphqls`, `gqlgen.yml`), **`/graphql`** Fiber **`adaptor`**, isteğe **GraphiQL** playground, **`graph-gophers/dataloader`** ile istek başına **`Loaders`**, **`zatrano gen graphql <model>`** |
| Operasyon | `/health`, `/ready`, `/status` |
| CLI | **`new`**, **`gen module`**, **`gen crud`**, **`gen request`**, **`gen policy`**, **`gen job`**, **`gen mail`**, **`gen event`**, **`gen listener`**, **`gen notification`**, **`gen model`**, **`gen middleware`**, **`gen resource`**, **`gen test`**, **`gen seeder`**, **`gen factory`**, **`gen command`**, **`gen graphql`**, `serve`, `db`, **`search import`**, **`cache`**, **`queue`**, **`mail`**, **`openapi export`**, `openapi validate`, **`jwt sign`**, **`api-key create`**, **`api-key list`**, **`api-key revoke`**, … |

**Şu an hazır:** `serve`, `doctor`, **`routes`**, **`config print`**, **`config validate`**, **`verify`** (isteğe **`--race`**), `completion`, `version` / **`--version`**, **`new`**, **`gen module`** + **`gen crud`** + **`gen request`** + **`gen policy`** + **`gen job`** + **`gen mail`** + **`gen event`** + **`gen listener`** + **`gen notification`** + **`gen model`** + **`gen middleware`** + **`gen resource`** + **`gen test`** + **`gen seeder`** + **`gen factory`** + **`gen command`** + **`gen wire`** + **`gen view`** + **`gen graphql`**, **`db`** (golang-migrate; varsayılan **embed** SQL **`pkg/migrations/sql/<sürücü>/`**, isteğe **file** + `migrations_dir` / `--migrations`) + **`db tenants`** (kiracı PostgreSQL şemasında migrate/rollback/create-schema), **`search import`** (Meilisearch/Typesense toplu indeks, `RegisterImporter`), **`pkg/features`** (bayraklar, rollout, şablon + HTTP middleware), **`pkg/graphql`** (gqlgen + dataloader kancaları), **`cache`** (Memory/Redis, Tags, middleware), **`queue`** (Redis FIFO, geciktirilmiş joblar, retry, başarısız joblar, worker), **`mail`** (SMTP/log, şablonlar, kuyruk, ek dosya, önizleme), **`events`** (senkron/asenkron gönderim, ShouldQueue, kuyruk tabanlı listener'lar), **`notifications`** (çok kanallı, Veritabanı/SMS/Push, okundu-takibi, Twilio/Netgsm/FCM/APNs), **`broadcast`** (WebSocket hub, Pusher tarzı protokol, private/presence JWT kanalları, SSE), **`audit`** (model activity + HTTP audit, JSON Patch farkları), **`pkg/search`** (PostgreSQL FTS scope'ları + harici sürücü), **`openapi validate`** + **`openapi export`**, **`jwt sign`**, **`storage`** (yerel/S3/MinIO/R2, imzalı URL'ler, resim işleme), **OAuth2**, **`http.*`** (CORS, rate limit, istek süresi, gövde boyutu), **`i18n`** (JSON yereller + Fiber yardımcıları), **validation** (generic `Validate[T]`, i18n hata mesajları, özel kurallar, form request'ler), **yetkilendirme** (RBAC rol→izin, Gate/Policy, `middleware.Can`, i18n 403), **çok kiracılık** (`middleware.ResolveTenant`, `pkg/tenant`, tenant kapsamlı repository), **view motoru** (`{{extends}}` layout kalıtımı, `{{block}}` bölümleri, `views/components/` partial'ları, form builder, flash mesajlar, eski girdi `{{old}}`, `{{asset}}` versiyonlanmış URL'ler, Vite/esbuild manifest + HMR), Redis + CSRF, JWT, Scalar **`/docs`**, **Air** (`.air.toml`).

---

## Dizilim (`pkg/` ve `internal/`)

| Yol | Amaç |
|-----|------|
| `pkg/config`, `pkg/core`, `pkg/server`, `pkg/health`, `pkg/middleware`, `pkg/security`, `pkg/auth`, `pkg/cache`, `pkg/queue`, `pkg/mail`, `pkg/notifications`, `pkg/events`, `pkg/broadcast`, `pkg/tenant`, `pkg/audit`, `pkg/search`, `pkg/features`, `pkg/graphql`, `pkg/oauth`, `pkg/openapi`, `pkg/i18n`, `pkg/validation`, `pkg/storage`, `pkg/database`, `pkg/migrations` (gömülü SQL; doğrudan import hedefi değil), `pkg/zatrano`, `pkg/meta` | **Genel API** — uygulamalar import eder |
| `internal/cli`, `internal/db`, `internal/gen` | **CLI ve üreticiler** — uygulama import etmez |

Üretilen projeler **`zatrano.Start`** + **`RegisterRoutes: routes.Register`** (`internal/routes/register.go`) veya ek rota yoksa **`zatrano.Run()`** kullanır.

---

## Gereksinimler

- Go **1.25.0+**
- **Bir veritabanı** — GORM ve `zatrano db migrate` için **PostgreSQL** (varsayılan), **MySQL**, **SQLite** veya **SQL Server** (`database_driver` + `database_url`; bkz. `pkg/database`, `config/examples/dev.yaml`)
- **Redis** — oturum + CSRF için (yerelde isteğe bağlı; prod'da genelde zorunlu)
- **PostgreSQL istemci araçları** — `zatrano db backup` ve `db restore` yalnızca Postgres yedeği için: `pg_dump`, `pg_restore`, `psql` PATH'te olmalı

---

## Kurulum

CLI'yi global olarak yükleyin:

```bash
go install github.com/zatrano/zatrano/cmd/zatrano@latest
```

---

## Hızlı başlangıç

Yeni uygulama oluşturun:

```bash
zatrano new app
cd app
zatrano serve
```

Veya framework'ü doğrudan çalıştırın:

```bash
go run ./cmd/zatrano serve
```

İsteğe bağlı:

```bash
cp config/examples/dev.yaml config/dev.yaml
cp .env.example .env
```

**`DATABASE_URL`** (ve isteğe **`DATABASE_DRIVER`**) ayarlandıktan sonra şemayı uygulayın (varsayılan **gömülü** migrasyonlar — `migrations_source: embed`):

```bash
zatrano db migrate --env dev --config-dir config
```

OpenAPI doğrulama ve dışa aktarma:

```bash
go run ./cmd/zatrano openapi validate api/openapi.yaml
go run ./cmd/zatrano openapi validate --merged
go run ./cmd/zatrano openapi export --output api/openapi.merged.yaml
```

---

## CLI komutları

| Komut | Açıklama |
|-------|----------|
| `zatrano serve` | HTTP sunucusu (`--addr`, `--env`, `--config-dir`, `--no-dotenv`) |
| `zatrano doctor` | Yapılandırma (**HTTP** ara katman özeti dahil) + bağlantı kontrolleri |
| `zatrano routes` | Rotalar (`serve` ile aynı config; `--json`, `--all`, **`--group`**) |
| `zatrano config print` | Maskeli tam çıktı; **`--paths-only`** kısa özet (varsayılan **satırlar**; `json` / `yaml`) |
| `zatrano config validate` | Yükle + **doğrula** (DB/Redis yok); CI için **`--quiet`** / **`-q`** (yalnızca çıkış kodu) |
| `zatrano new <name>` | Yeni uygulama (`--module`, `--output`, yerel geliştirme için `--replace-zatrano`) |
| `zatrano db migrate` | Varsayılan olarak `pkg/migrations/sql/<sürücü>/` içindeki gömülü SQL migration’ları uygula (`migrations_source: embed`); disk için `file` + `migrations_dir` veya `--migrations` |
| `zatrano db rollback` | Geri al (`--steps`) |
| `zatrano db seed` | `db/seeds/*.sql` (yoksa no-op) |
| `zatrano db backup` | `pg_dump` → dosya (`--format`, `--output` veya varsayılan `backups/`) |
| `zatrano db restore` | `pg_restore` / `psql` (**`--yes` zorunlu**, isteğe `--clean`) |
| `zatrano db tenants migrate` | PostgreSQL **`search_path`** ile kiracı şemasında migration (`--tenant` zorunlu; `--env`, `--config-dir`, `--migrations`, `--steps` ile `db migrate` ile aynı mantık) |
| `zatrano db tenants rollback` | Kiracı şemasında migration geri alma (`--tenant`, `--steps`, …) |
| `zatrano db tenants create-schema` | Hesaplanan kiracı şeması için `CREATE SCHEMA IF NOT EXISTS` |
| `zatrano search import <model>` | `search.RegisterImporter` ile kayıtlı modeli Meilisearch veya Typesense'e toplu yazar (`search.enabled`, `database_url`; `--env`, `--config-dir`, `--no-dotenv`) |
| `zatrano gen module <name>` | `modules/<name>/` + **wire** + wire dosyasında **`go fmt`** (`--skip-wire`, `--module-root`, `--out`, `--dry-run`) |
| `zatrano gen crud <name>` | CRUD + **form request struct'ları** (`requests/`) + **`RegisterCRUD`** wire + **`go fmt`** (aynı bayraklar) |
| `zatrano gen request <name>` | Yalnızca form request struct'ları üret (`modules/<name>/requests/create_*.go`, `update_*.go`) |
| `zatrano gen policy <name>` | Yetkilendirme policy stub'ı üret (`modules/<name>/policies/<name>_policy.go`) — `auth.Policy` arayüzünü CRUD metotlarıyla implemente eder |
| `zatrano gen job <name>` | Kuyruk job stub'ı üret (`modules/jobs/<name>.go`) — `queue.Job` arayüzünü Handle, Retries, Timeout ile implemente eder |
| `zatrano gen mail <name>` | Mailable struct + HTML şablon üret (`modules/mails/<name>_mail.go` + `views/mails/<name>.html`) |
| `zatrano gen event <name>` | Event struct üret (`modules/events/<name>_event.go`) — `events.Event` arayüzünü implemente eder |
| `zatrano gen listener <name>` | Listener üret (`modules/listeners/<name>_listener.go`); asenkron için `--queued` bayrağı |
| `zatrano gen notification <name>` | Çok kanallı gönderim için bildirim stub'ı üret (`modules/notifications/<name>.go`) |
| `zatrano gen model <name>` | Model scaffold ve PostgreSQL migration iskeleti üret (`pkg/repository/models/`, `pkg/migrations/sql/postgres/`) |
| `zatrano gen middleware <name>` | Fiber middleware stub'ı üret (`pkg/middleware/`) |
| `zatrano gen resource <name>` | API resource transformer stub'ı üret (`pkg/resources/`) |
| `zatrano gen test <name>` | Handler ve servis test stub'ları üret (`tests/`) |
| `zatrano gen seeder <name>` | SQL seed dosyası üret (`db/seeds/`) |
| `zatrano gen factory <name>` | Test veri factory stub'ı üret (`pkg/factory/`) |
| `zatrano gen command <name>` | Cobra CLI komut scaffold'u üret (`internal/cli/`) |
| `zatrano gen wire <name>` | Sadece wire (dosya üretmez); `register.go` / `crud_register.go` varlığına göre (`--register-only`, `--crud-only`) |
| `zatrano gen graphql <model>` | `api/graphql/<model>_stub.graphqls` + **`go run github.com/99designs/gqlgen@v0.17.78 generate`** (`--module-root`, `--dry-run`, `--skip-generate`) |
| `zatrano gen view <n>` | `views/<n>/` altında sunucu taraflı HTML şablonları scaffold oluştur (`index.html`, `show.html`; `--with-form` ile `create.html` + `edit.html`; `--layout`, `--dry-run`) |
| `zatrano openapi validate` | Tek dosya veya **`--merged`** (canlı `/openapi.yaml` ile aynı; `--base`, isteğe konumsal argüman) |
| `zatrano openapi export` | Birleşik YAML yaz (`--base`, `--output` veya `-` stdout) |
| `zatrano jwt sign` | Test JWT üret (`--sub`, `--secret`, config bayrakları) |
| `zatrano cache clear` | Önbelleği temizle veya belirli tag'leri sil (`--tag`) |
| `zatrano queue work` | Kuyruk worker süreci başlat (`--queue`, `--tries`, `--timeout`, `--sleep`) |
| `zatrano queue failed` | Başarısız jobları listele |
| `zatrano schedule run` | Kayıtlı planlanmış taskları çalıştır |
| `zatrano schedule list` | Kayıtlı planlanmış taskları ve cron ifadelerini listele |
| `zatrano queue retry [id]` | Başarısız jobı yeniden gönder veya `--all` |
| `zatrano queue flush` | Tüm başarısız job kayıtlarını sil |
| `zatrano mail preview [name]` | E-posta şablonunu tarayıcıda önizle (`--port`, `--layout`) |
| `zatrano storage:link` | `storage/app/public` dosyasından `public/storage` sembolik bağlantı oluştur (`--force`, `--storage-path`, `--public-path`) |
| `zatrano storage:clear [disk]` | Depolama diskinden tüm dosyaları sil (`--force` ile onay iste) |
| `zatrano completion …` | Kabuk tamamlama |
| `zatrano verify` | **`go vet` + `go test` + birleşik OpenAPI** (PR/CI; yarış için **`--race`**; `--no-vet`, `--no-test`, `--no-openapi`, `--module-root`) |
| `zatrano version` | Sürüm (ayrıca **`zatrano --version`**) |

**Windows / boşluklu yol:** `--replace-zatrano` ile framework kökünü verin; gerekirse `go.mod` içinde yol tırnaklanır.

---

## HTTP (şimdilik)

| Metot | Yol | Not |
|-------|-----|-----|
| GET | `/` | JSON özet (`env`, `endpoints`, `http` CORS/rate-limit bayrakları, `error_includes_request_id`) |
| GET | `/health`, `/ready`, `/status` | Canlılık / hazırlık / özet (`/status` içinde `env`) |
| GET | `/openapi.yaml`, `/docs` | **Birleşik** OpenAPI (`/` ve `/status` için JSON şema) + Scalar |
| GET | `/api/v1/public/ping` | Herkese açık |
| GET | `/api/v1/private/me` | `jwt_secret` varsa **Bearer JWT** |
| POST | `/api/v1/auth/token` | Yalnızca `security.demo_token_endpoint: true` ve **`env: prod` değil** |
| GET | `/auth/oauth/google/login`, `/auth/oauth/github/login` | OAuth2 başlatır (`oauth.enabled` + anahtarlar gerekli) |
| GET | `/auth/oauth/google/callback`, `/auth/oauth/github/callback` | OAuth yönlendirme |
| GET | `/broadcast/ws` | **WebSocket** (`broadcast.enabled: true` iken); Pusher tarzı JSON; JWT sorgu `access_token` veya `Authorization` |
| GET | `/broadcast/sse/:kanal` | **SSE** (`broadcast.enabled` + `broadcast.sse_enabled`); WebSocket ile aynı kanal adları; token sorgu veya başlık |

**Oturum + CSRF:** `redis_url` ve `security` uygunsa açılır. CSRF, `Bearer`, `csrf_skip_prefixes` (varsayılan `/api/`) ve **`/auth/oauth/`** için atlanır.

**OAuth2:** `oauth.enabled`, `oauth.base_url`, `oauth.providers.google` / `github` ayarlayın. Sağlayıcı konsolunda yönlendirme: `{base_url}/auth/oauth/google/callback` (GitHub için aynı kalıp). Oturum alanları: `oauth_provider`, `oauth_subject`, `oauth_name`, `oauth_email`.

**Hatalar:** JSON gövdesi `{ "error": { "code", "message", "request_id"? } }`. `request_id`, **`X-Request-ID`** başlığıyla aynıdır (log ve destek için).

**HTTP ara katmanı (`http` YAML / `HTTP_*` env):**

- **CORS** — `http.cors_enabled`, `cors_allow_origins`, `cors_allow_methods`, `cors_allow_headers`, `cors_expose_headers`, `cors_allow_credentials`, `cors_max_age`. Varsayılan **kapalı**. **`cors_allow_credentials: true`** ile köken **`*`** birlikte kullanılamaz (doğrulama hata verir).
- **Rate limit** — `rate_limit_enabled`, `rate_limit_max`, `rate_limit_window`, isteğe **`rate_limit_redis: true`** (`redis_url` gerekir). Aksi halde süreç başına **bellek içi**. Limit **altındaki** yanıtlarda **`X-RateLimit-*`** vardır. Limit aşımında **429** + aynı `error` JSON + **`Retry-After`** (RFC 6585).
- **İstek süresi** — `request_timeout` (ör. `60s`): Fiber **timeout**; aşımda **408** JSON.
- **Gövde boyutu** — `body_limit` bayt (`0` = Fiber varsayılanı **4 MiB**).
- **Zarif HTTP kapatma** — `http.shutdown_timeout` (varsayılan `15s`): Fiber `ShutdownWithContext` üst sınırı. Gömülü sunucu için `zatrano.StartOptions.ShutdownHooks` ile aynı süre içinde ek adımlar çalıştırılabilir.
- **Sıfır kesinti yeniden başlatma (Unix)** — `http.graceful_restart: true`: [tableflip](https://github.com/cloudflare/tableflip) ile dinleyici soketinin yeni sürece devri; tetiklemek için sürece **`SIGUSR2`** gönderin (yeni ikili başlar, eski süreç bağlantıları boşaltır). İsteğe `http.graceful_restart_pid_file` (systemd `PIDFile=` senaryosu). **`go run` desteklenmez**; gerçek derlenmiş ikili kullanın. Windows yapılandırmada açık olsa bile çalışma zamanında devre dışı bırakılır.

Yığın sırası: **recover → request-id → i18n (açıksa) → CORS → timeout → rate limit → helmet → compress → oturum/CSRF → rotalar**.

---

## Broadcasting / WebSocket

ZATRANO, **`pkg/broadcast`** altında **bellek içi (in-process) bir yayın hub'ı** sunar: aynı kanal adına hem **WebSocket** hem isteğe bağlı **SSE** aboneleri bağlanır. Kablo formatı **Pusher uyumlu bir alt küme**dir; **Laravel Echo**, **pusher-js** veya aynı JSON olaylarını konuşan istemciler kullanılabilir.

### Açma

```yaml
# config/dev.yaml
broadcast:
  enabled: true
  path_prefix: /broadcast
  jwt_query_param: access_token
  sse_enabled: true
  allow_origins: []
```

`broadcast.enabled: true` olduğunda bootstrap `app.Broadcast` oluşturur ve rotaları kaydeder. **private-** / **presence-** kanalları için **`jwt_secret`** tanımlı olmalıdır (HS256, `security.JWTMiddleware` ile aynı kurallar).

### Kanal adları

| Önek | Kimlik doğrulama | Not |
|------|------------------|-----|
| *(yok)* / genel isimler | yok | Herkes abone olabilir. |
| `private-…` | Bağlantıda geçerli JWT | İsteğe `private-user-{sub}` kalıbı: yalnızca o JWT `sub` değeri. |
| `presence-…` | JWT | Aynı `presence-user-{sub}` kalıbı; `pusher_internal:subscription_succeeded` içinde üye listesi; **`member_added` / `member_removed`** olayları. |

### Sunucudan yayın

```go
// app.Broadcast, broadcast.enabled açıkken *broadcast.Hub
_ = app.Broadcast.PublishJSON("public-news", "ArticlePublished", map[string]any{"id": 42})
```

### Presence yardımcısı

```go
ids := app.Broadcast.OnlineOn("presence-room")
```

O presence kanalındaki güncel **JWT `sub`** listesini döner (yalnızca tek süreç içi; ölçekte Redis vb. gerekir).

### WebSocket protokolü

Sunucu önce **`pusher:connection_established`** ve **`socket_id`** gönderir. Abonelik:

```json
{ "event": "pusher:subscribe", "data": { "channel": "public-news" } }
```

Presence için **`channel_data`**: `{"user_id":"…","user_info":{…}}` (opsiyonel `user_id`, doluysa JWT `sub` ile eşleşmeli).

**Ping:** `{"event":"pusher:ping","data":{}}` → **`pusher:pong`**.

### SSE

`GET /broadcast/sse/kanal-adı` aynı JSON zarflarını **`data:`** satırları olarak akıtır. **private/presence** için **`?access_token=`** (veya `jwt_query_param`) kullanın.

---

## Multi-tenancy

Her HTTP isteğinde **kiracı çözümleme**, isteğe bağlı **satır düzeyinde** (`tenant_id` vb.) repository filtresi ve **PostgreSQL şema başına** migration desteği sunulur.

### Yapılandırma

```yaml
tenant:
  enabled: true
  mode: header              # header | subdomain
  header_name: X-Tenant-ID
  subdomain_suffix: ".app.local"   # mode=subdomain iken zorunlu
  required: false
  isolation: row            # row | schema
  row_column: tenant_id
  schema_prefix: tenant_
```

**`tenant.enabled: true`** iken **`middleware.ResolveTenant`**, **request-id** sonrasında çalışır. Kiracı bilgisi **`c.Locals(middleware.LocalsTenant)`** ve **`tenant.WithContext`** ile **`c.Context()`** üzerinde taşınır (GORM/repository ile aynı `context.Context` kullanın).

### Satır izolasyonu (ortak şema)

**`repository.NewTenantAware[T](db, "tenant_id")`** kullanın; tüm okuma/yazma sorgularına çözülen kiracıya göre **`WHERE`** eklenir. Sayısal olmayan anahtarlar metin sütununda eşitlik için kullanılır (`tenant_slug` gibi).

**`repository.TenantFK`** gömerek **`tenant_id`** alanını sayısal kiracı anahtarından **BeforeCreate** ile doldurun.

### Şema izolasyonu

1. **`isolation: schema`** ile **`Info.Schema`** (ör. `tenant_acme`) üretilir.
2. İstek başına: **`tenant.GormSession(app.DB, c.Context())`** ile `search_path` ayarlanır; **`repository.New`** bu oturumla kullanılabilir.
3. Şema ve migration:

```bash
zatrano db tenants create-schema --tenant acme
zatrano db tenants migrate --tenant acme
```

DSN'e **`options=-csearch_path=<şema>,public`** eklenir; **golang-migrate** DDL'i kiracı şemasında çalıştırır (varsayılan **embed** kaynak: `pkg/migrations/sql/postgres/`; `--migrations` ile dosya kaynağına geçilebilir).

**Not:** Ölçekte RLS, ayrı veritabanı veya paylaşımlı önbellek stratejisi uygulama sorumluluğundadır.

---

## Audit / Activity log

**Model activity** (kim / ne zaman / ne değişti) ve isteğe bağlı **HTTP istek audit** kaydı sunulur.

### Yapılandırma

```yaml
audit:
  enabled: true
  model_enabled: true
  http_enabled: true
  http_driver: db          # db | file
  http_file_path: logs/http_audit.jsonl
```

**`zatrano db migrate`** ile **`000009_audit`** migration'ı **`zatrano_activity_logs`** ve **`zatrano_http_audit_logs`** tablolarını oluşturur.

### Model activity

1. Başlangıçta türleri kaydedin: `audit.RegisterSubject[Product]("products")`.
2. GORM'a verdiğiniz `context` içine `audit.WithUser` / `audit.WithRequest` ekleyin. `middleware.AuditLog` açıksa istek bağlamına **request id**, **IP** ve varsa **JWT `sub`** eklenir.
3. **changes** alanı **RFC 6902 JSON Patch** (yüzeysel nesne anahtarları). **`audit.DiffJSONPatch`** ile kendi karşılaştırmalarınızı üretebilirsiniz.
4. **`audit.Skip(ctx)`** ile tek zincirde callback'leri atlayın.

**Soft delete** (`deleted_at` güncellemesi) çoğu durumda **update** olarak loglanır.

### HTTP audit

**`middleware.AuditLog`**, `audit.enabled` + `audit.http_enabled` iken **`AccessLog`** sonrasında çalışır; **DB** veya **JSONL dosya** yazar. Kullanıcı: önce **`LocalsUserID`**, yoksa JWT **`sub`**.

### Diff

Bkz. **`audit.DiffJSONPatch`** (`pkg/audit/patch.go`).

---

## Full-text search

**PostgreSQL** tarafında `tsvector` sütununuzu (ör. `GENERATED ALWAYS AS (to_tsvector(...)) STORED` veya tetikleyici) migration ile tanımlarsınız; sorgu tarafında **`pkg/search`** GORM **`repository.Scope`** üreticileri kullanılır:

- **`search.WhereFullTextMatch(regconfig, vectorColumn, kullanıcıSorgusu)`** — boş sorguda no-op; aksi halde `vectorColumn @@ plainto_tsquery(regconfig, kullanıcıSorgusu)`.
- **`search.OrderByTSRank(regconfig, vectorColumn, kullanıcıSorgusu)`** — `ts_rank_cd` ile azalan sıra (sorgu boşsa no-op).

`vectorColumn` değeri **güvenilir bir SQL tanımlayıcısı** olmalıdır (sabit veya whitelist; kullanıcı girdisi olarak vermeyin). `regconfig` ve metin **`plainto_tsquery`** ile parametre bağlanır.

Uygulama genelinde varsayılan dil için **`search.postgres_fts_language`** (`config.Search`, örn. `simple`, `english`, `turkish`) kullanılabilir.

### Meilisearch / Typesense

```yaml
search:
  enabled: true
  driver: meilisearch          # veya typesense
  default_index_prefix: zatrano_
  meilisearch_url: http://127.0.0.1:7700
  meilisearch_api_key: ""
  # typesense_url / typesense_api_key — driver typesense iken zorunlu
  postgres_fts_language: simple
```

**`core.Bootstrap`** açıkken **`app.Search`** ( **`search.NewClient`** ) harici sürücüyü sarar; indeks UID / koleksiyon adı = `default_index_prefix` + mantıksal ad.

Toplu aktarım için uygulama **`init`** içinde **`search.RegisterImporter("product", func(ctx, db, drv) error { ... })`** tanımlar; CLI:

```bash
zatrano search import product
```

İçerik motorunda hedef indeks/koleksiyon ve şema (ör. Meilisearch'te birincil anahtar **`id`**) önceden oluşturulmuş olmalıdır.

---

## Feature flags

**`pkg/features`** ile bayrakları **YAML** (`features.definitions`), **PostgreSQL** (`zatrano_feature_flags`, migration **`000010_feature_flags`**) veya **ikisi birden** (`source: both` — satır varsa DB kazanır) yönetebilirsiniz. **`features.enabled: false`** iken tüm bayraklar kapalıdır; **`app.Features`** yine oluşturulur (no-op).

### Yapılandırma

```yaml
features:
  enabled: true
  source: both              # config | db | both
  definitions:
    - key: beta-ui
      enabled: true
    - key: new-dashboard
      enabled: true
      rollout_percent: 30   # 1..99: yalnızca giriş yapmış kullanıcı + kararlı FNV kovası (A/B)
      allowed_roles: [admin, editor]
```

- **`allowed_roles`** doluysa **anonim** isteklerde bayrak **kapalı**dır; roller **`middleware.InjectRoles`** (veya aynı Locals anahtarlarını dolduran akış) ile gelmelidir.
- **`rollout_percent`** 1–99 arasıyken **kullanıcı id’si** yoksa sonuç **kapalı**dır (kovaya girmek için `uint` id gerekir).

### Go API

```go
u := &features.User{ID: 1, Roles: []string{"admin"}}
if app.Features.For(u).IsEnabled(c.Context(), "new-dashboard") {
    // ...
}
// İstekten: app.Features.FromFiber(c).IsEnabled(ctx, "beta-ui")
```

### HTTP middleware

```go
import "github.com/zatrano/zatrano/pkg/middleware"

app.Get("/beta", middleware.RequireFeature(app.Features, "beta-ui"), handler)
```

Bayrak kapalıysa **404** döner (rota “yok” gibi). **`server.Mount`**, `features.enabled` iken **`Features.LocalsMiddleware`** ekler; **`features.EvalFromFiber(c)`** ile aynı değerlendiriciye erişebilirsiniz.

### View şablonu

**`a.View.ViewData(c, ...)`** kullanın; içine otomatik **`feature`** şablon fonksiyonu için gerekli bağlama eklenir:

```html
{{if feature . "beta-ui"}}
  <p>Beta arayüzü açık</p>
{{end}}
```

İlk argüman **şablon kökü (`.`)** olmalıdır; `html/template` kısıtı nedeniyle `{{if feature "beta-ui"}}` tek başına kullanılamaz.

---

## GraphQL

**Şema öncelikli** GraphQL, kökte **`gqlgen.yml`** ve **`api/graphql/*.graphqls`** dosyaları ile tanımlanır; üretilen sunucu kodu **`pkg/graphql/graph/`** altındadır. Şema değişikliğinden sonra:

```bash
go run github.com/99designs/gqlgen@v0.17.78 generate
```

(İlk **`go run github.com/99designs/gqlgen@…`** indirme ve çözümlemesi biraz sürebilir; sonraki çağrılar önbellekten hızlıdır.)

### Yapılandırma

```yaml
graphql:
  enabled: true
  path: /graphql
  playground: true        # prod’da genelde false
  playground_path: /playground
```

**`server.Mount`**, `graphql.enabled` iken **`graphql.Register`** çağrısıyla Fiber üzerinde **`middleware/adaptor`** ile `net/http` gqlgen işleyicisini bağlar; kök JSON’da **`graphql`** ve isteğe **`graphql_playground`** uçları listelenir.

### DataLoader

Her HTTP isteği için **`graphql.NewLoaders(app)`** oluşturulur ve **`graphql.WithLoaders`** ile `context.Context` üzerinden resolver’lara taşınır. Resolver içinde:

```go
import zgraphql "github.com/zatrano/zatrano/pkg/graphql"

ld := zgraphql.LoadersFrom(ctx)
_ = ld // alanlarınız: *dataloader.Loader[...] (graph-gophers/dataloader/v7)
```

### Kod üretimi

```bash
zatrano gen graphql product --module-root .
```

Bu komut **`api/graphql/product_stub.graphqls`** ekler (ör. `extend type Query { product(id: ID!): Product }` ve `type Product { id: ID! }`), ardından **gqlgen generate** çalıştırır. **`--skip-generate`** yalnızca `.graphqls` yazar. Aynı isimde stub zaten varsa hata verir.

---

## Validation (Doğrulama)

ZATRANO, [`go-playground/validator/v10`](https://pkg.go.dev/github.com/go-playground/validator/v10) kütüphanesini sarmalayan **generic, struct-tag tabanlı bir doğrulama sistemi** sunar. Otomatik **422 JSON yanıtları** ve **i18n çeviri desteği** içerir.

### Temel Kullanım

Ana API **`zatrano.Validate[T](c)`** — tek bir generic çağrı ile istek gövdesi parse edilir, struct tag'leri doğrulanır ve hata durumunda yapılandırılmış 422 yanıtı döner:

```go
import "github.com/zatrano/zatrano/pkg/zatrano"

func (h *ProductHandler) Create(c fiber.Ctx) error {
    req, err := zatrano.Validate[CreateProductRequest](c)
    if err != nil {
        return err // 422 JSON yanıtı zaten gönderildi
    }
    // req geçerli — kullanabilirsiniz
    return h.svc.Create(c.Context(), req.Name, req.Email)
}
```

### Form Request Struct'ları

İstek yapılarınızı `json` ve `validate` tag'leri ile tanımlayın:

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

### Request Struct'larını Üretme

CLI ile request stub'larını otomatik oluşturabilirsiniz:

```bash
# Sadece request struct'larını üret
zatrano gen request product
# → modules/product/requests/create_product.go
# → modules/product/requests/update_product.go

# gen crud komutu request struct'larını da otomatik üretir
zatrano gen crud product
# → modules/product/crud_handlers.go      (zatrano.Validate[T] kullanır)
# → modules/product/crud_register.go
# → modules/product/requests/create_product.go
# → modules/product/requests/update_product.go
```

### 422 Hata Yanıt Formatı

Doğrulama başarısız olduğunda, tutarlı bir JSON yapısı döner:

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
        "value": "geçersiz-email",
        "message": "Must be a valid email address"
      }
    ]
  }
}
```

i18n açık ve istek dili `tr` olduğunda, mesajlar otomatik çevrilir:

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
        "value": "geçersiz-email",
        "message": "Geçerli bir e-posta adresi olmalıdır"
      }
    ]
  }
}
```

### i18n Doğrulama Mesajları

Doğrulama mesajları, locale dosyalarında `validation.*` anahtar alanı altında tutulur:

```json
// locales/tr.json
{
  "validation": {
    "required": "Bu alan zorunludur",
    "email": "Geçerli bir e-posta adresi olmalıdır",
    "min": "En az {{.Param}} karakter olmalıdır",
    "max": "En fazla {{.Param}} karakter olmalıdır"
  }
}
```

`{{.Param}}` yer tutucusu kısıtlamanın değeri ile değiştirilir (ör. `min=5` → `"En az 5 karakter olmalıdır"`).

**Hazır çevrilmiş tag'ler:** `required`, `email`, `min`, `max`, `gte`, `lte`, `gt`, `lt`, `len`, `url`, `uri`, `uuid`, `oneof`, `numeric`, `number`, `alpha`, `alphanum`, `boolean`, `contains`, `excludes`, `startswith`, `endswith`, `ip`, `ipv4`, `ipv6`, `datetime`, `json`, `jwt`, `eqfield`, `nefield`.

### Özel Doğrulama Kuralları

İsteğe bağlı i18n desteği ile özel doğrulama tag'leri kaydedin:

```go
import (
    "github.com/go-playground/validator/v10"
    "github.com/zatrano/zatrano/pkg/zatrano"
)

// Özel kural kaydet
zatrano.RegisterRule("tc_no", func(fl validator.FieldLevel) bool {
    v := fl.Field().String()
    if len(v) != 11 {
        return false
    }
    // ... TC kimlik numarası algoritması
    return true
})

// i18n mesaj anahtarı ile (locale dosyalarınıza "validation.tc_no" ekleyin)
zatrano.RegisterRuleWithMessage("tc_no", tcNoValidator, "validation.tc_no")
```

Sonra struct tag'lerinde kullanın:

```go
type VatandasRequest struct {
    TCNO string `json:"tc_no" validate:"required,tc_no"`
}
```

### Doğrudan Engine Erişimi

İleri düzey kullanım için alttaki validator engine'e erişebilirsiniz:

```go
import "github.com/zatrano/zatrano/pkg/validation"

engine := validation.Default()
engine.Validator() // go-playground/validator'dan *validator.Validate

// Herhangi bir struct'ı programatik olarak doğrulayın (Fiber context'i olmadan)
if verr := engine.ValidateStruct(myStruct, "tr"); verr != nil {
    for _, fe := range verr.Errors {
        fmt.Printf("%s: %s\n", fe.Field, fe.Message)
    }
}
```

---

## Yetkilendirme (RBAC & Gate/Policy)

ZATRANO, iki tamamlayıcı katmanlı **eksiksiz bir yetkilendirme sistemi** sunar: izin kontrolleri için **RBAC** (rol tabanlı, DB destekli) ve kaynak düzeyinde ince taneli yetkilendirme için **Gate/Policy** (kaynak bazlı). Her ikisi de yerelleştirilmiş 403 hata mesajları için **i18n** sistemiyle entegredir.

### RBAC — Rol Tabanlı Erişim Kontrolü

Roller ve izinler veritabanında saklanır (`roles`, `permissions`, `role_permissions`, `zatrano_user_roles` tabloları). Yoğun yolda DB çağrılarından kaçınmak için bellek içi önbellek kullanılır. `RBACManager`, bootstrap sırasında otomatik olarak başlatılır (DB varken) ve `app.RBAC` ile erişilebilir.

```go
import "github.com/zatrano/zatrano/pkg/auth"

// Rol ve izin oluşturma
rbac := app.RBAC
rbac.CreateRole(ctx, "admin", "Yönetici")
rbac.CreateRole(ctx, "editor", "İçerik editörü")
rbac.CreatePermission(ctx, "posts.create", "Yazı oluştur")
rbac.CreatePermission(ctx, "posts.update", "Yazı güncelle")
rbac.CreatePermission(ctx, "posts.delete", "Yazı sil")

// Rollere izin atama
rbac.AssignPermissions(ctx, "admin", "posts.create", "posts.update", "posts.delete")
rbac.AssignPermissions(ctx, "editor", "posts.create", "posts.update")

// Kullanıcıya rol atama
rbac.AssignRoleToUser(ctx, userID, "editor")

// İzin kontrolü
ok, _ := rbac.UserHasPermission(ctx, userID, "posts.create") // true
ok, _ = rbac.UserHasPermission(ctx, userID, "posts.delete")  // false (editör silemez)
```

**Veritabanı migration:** `zatrano db migrate` çalıştırın — `000002_rbac` migration'ı dört tabloyu uygun indeks ve foreign key'lerle oluşturur.

### Gate / Policy — Kaynak Bazlı Yetkilendirme

`Gate` sistemi (`app.Gate` ile erişilir) belirli eylemler için yetkilendirme kontrolleri tanımlamaya olanak sağlar. Anlık kontroller için `Define`, yapılandırılmış CRUD policy'leri için `RegisterPolicy` kullanın.

```go
import "github.com/zatrano/zatrano/pkg/auth"

// Anlık gate tanımı
gate := app.Gate
gate.Define("edit-post", func(c fiber.Ctx, resource any) bool {
    post := resource.(*Post)
    userID, _ := c.Locals(middleware.LocalsUserID).(uint)
    return post.AuthorID == userID
})

// Süper-admin bypass (her gate kontrolünden önce çalışır)
gate.Before(func(c fiber.Ctx, ability string, resource any) *bool {
    roles, _ := c.Locals(middleware.LocalsUserRoles).([]string)
    for _, r := range roles {
        if r == "super-admin" { t := true; return &t }
    }
    return nil // gate tanımına düş
})

// Handler'larda:
if err := gate.Authorize(c, "edit-post", post); err != nil {
    return err // 403 Forbidden
}
```

**Policy arayüzü:** yapılandırılmış CRUD yetkilendirmesi için `auth.Policy` arayüzünü implemente edin. `zatrano gen policy` ile stub üretin:

```bash
zatrano gen policy post
# → modules/post/policies/post_policy.go
```

Üretilen policy 7 metot içerir: `ViewAny`, `View`, `Create`, `Update`, `Delete`, `ForceDelete`, `Restore`. Gate'e kaydedin:

```go
import "myapp/modules/post/policies"

gate.RegisterPolicy("post", &policies.PostPolicy{})
// Oluşturur: "post.viewAny", "post.view", "post.create", "post.update",
//            "post.delete", "post.forceDelete", "post.restore"
```

### Route Seviyesi Yetkilendirme Middleware

`pkg/middleware` paketi, rota seviyesinde izin ve rol kontrolü için kullanıma hazır middleware sağlar. Tümü **403 JSON** + **i18n destekli** hata mesajı döner.

| Middleware | Açıklama |
|---|---|
| `middleware.Can(rbac, "perm")` | Kullanıcının belirli bir izne sahip olmasını gerektirir |
| `middleware.CanAny(rbac, "p1", "p2")` | Kullanıcı listelenen izinlerden **herhangi birine** sahipse geçer |
| `middleware.CanAll(rbac, "p1", "p2")` | Kullanıcı **tüm** listelenen izinlere sahipse geçer |
| `middleware.HasRole("admin")` | Belirli bir rol gerektirir |
| `middleware.HasAnyRole("admin", "editor")` | Kullanıcı listelenen rollerden **herhangi birine** sahipse geçer |
| `middleware.GateAllows(gate, "ability")` | Gate ability kontrolü (kaynak olmadan) |
| `middleware.InjectRoles(rbac)` | Kullanıcı rollerini Locals'a yükler (auth middleware'den sonra yerleştirin) |

```go
import "github.com/zatrano/zatrano/pkg/middleware"

// Kimlik doğrulama middleware'inden sonra:
app.Use(security.JWTMiddleware(cfg))
app.Use(middleware.InjectRoles(rbac))  // rolleri context'e yükler

// İzin bazlı
app.Get("/admin/users", middleware.Can(rbac, "users.view"), usersHandler)
app.Post("/posts", middleware.Can(rbac, "posts.create"), createPostHandler)
app.Delete("/system", middleware.CanAll(rbac, "system.admin", "system.delete"), handler)

// Rol bazlı
app.Get("/dashboard", middleware.HasAnyRole("admin", "editor"), dashHandler)

// Gate bazlı
app.Get("/posts", middleware.GateAllows(gate, "post.viewAny"), listPostsHandler)
```

### 403 Hata Yanıt Formatı

Yetkilendirme başarısız olduğunda, standart JSON hata yapısı döner:

```json
{
  "error": {
    "code": 403,
    "message": "You do not have permission to perform this action.",
    "permission": "posts.delete"
  }
}
```

i18n açık ve istek dili `tr` olduğunda:

```json
{
  "error": {
    "code": 403,
    "message": "Bu işlemi gerçekleştirme yetkiniz bulunmamaktadır.",
    "permission": "posts.delete"
  }
}
```

### i18n Yetkilendirme Mesajları

Yetkilendirme mesajları, locale dosyalarında `auth.*` anahtar alanı altında tutulur:

```json
{
  "auth": {
    "forbidden": "Bu işlemi gerçekleştirme yetkiniz bulunmamaktadır.",
    "unauthorized": "Bu kaynağa erişmek için kimlik doğrulaması gereklidir.",
    "role_required": "Bu kaynağa erişmek için gerekli role sahip değilsiniz.",
    "permission_required": "Gerekli izne sahip değilsiniz: {{.Permission}}."
  }
}
```

---

## Cache Sistemi (Önbellek)

ZATRANO, **Bellek İçi (In-Memory)** ve **Redis** sürücüleri için ortak bir API sunan **güçlü bir önbellek katmanı** sağlar. `Remember`, JSON serileştirme, tag tabanlı geçersiz kılma ve response middleware gibi gelişmiş özellikleri destekler.

### Sürücüler (Drivers)

Sistem, yapılandırmanıza göre en iyi sürücüyü otomatik olarak seçer:
- **Redis:** `redis_url` ayarlandığında tercih edilir. Dağıtık ortamlar ve tag desteği için uygundur.
- **Memory:** Yerel geliştirme veya tek node'lu yapılar için geri dönüş (fallback) seçeneğidir. Hızlıdır ancak geçicidir.

### Temel Kullanım

Önbellek yöneticisine `app.Cache` üzerinden erişebilirsiniz:

```go
import "context"

ctx := context.Background()

// Basit veri saklama
app.Cache.Set(ctx, "anahtar", "değer", 10 * time.Minute)

// Veri okuma
val, ok := app.Cache.Get(ctx, "anahtar")

// Otomatik JSON işleme
type Kullanici struct { Ad string }
app.Cache.SetJSON(ctx, "user:1", Kullanici{Ad: "Deniz"}, time.Hour)

var user Kullanici
ok, err := app.Cache.GetJSON(ctx, "user:1", &user)
```

### Gelişmiş Kullanım

#### `Remember` ve `RememberJSON`

Laravel stili popüler desen: Veri önbellekte varsa döner, yoksa verilen fonksiyonu çalıştırıp sonucu önbelleğe yazar ve döner.

```go
// Veri yoksa DB'den çek ve önbelleğe yaz
var users []User
err := app.Cache.RememberJSON(ctx, "users:all", 30*time.Minute, &users, func() (any, error) {
    return db.FindAllUsers(ctx)
})
```

#### Tags (Sadece Redis)

İlişkili anahtarları tag'ler altında gruplayarak toplu silme yapmanızı sağlar.

```go
// Tag ile saklama
app.Cache.Tags("users").Set(ctx, "users:1", data, time.Hour)

// Bir tag'e ait tüm anahtarları temizleme
app.Cache.Tags("users").Flush(ctx)
```

### Middleware (Ara Katman)

Bir rotanın tüm HTTP yanıtını sunucu tarafında önbelleğe alabilirsiniz.

```go
import "github.com/zatrano/zatrano/pkg/middleware"

// 5 dakika boyunca önbelleğe al
app.Get("/api/v1/stats", middleware.Cache(app.Cache, 5*time.Minute), handler)

// Tag desteği ile
app.Get("/api/v1/users", middleware.CacheWithConfig(app.Cache, middleware.CacheConfig{
    TTL:  10 * time.Minute,
    Tags: []string{"users"},
}), handler)
```

### CLI Komutları

Terminal üzerinden önbelleği temizleyin:

```bash
# Tüm önbelleği sil
zatrano cache clear

# Sadece belirli tag'leri sil
zatrano cache clear --tag users --tag posts
```

---

## Kuyruk / Job Sistemi

ZATRANO, geciktirilmiş zamanlama, otomatik yeniden deneme ve üssel geri çekilme (exponential backoff) ve başarısız job’ların PostgreSQL’de saklanmasıyla **Redis tabanlı bir arkaplan job kuyruğu** sunar.

### Job Tanımlama

`queue.Job` arayüzünü implemente edin veya varsayılan değerler için `queue.BaseJob` gömün:

```go
package jobs

import (
    "context"
    "time"
    "github.com/zatrano/zatrano/pkg/queue"
)

type EpostaGonderJob struct {
    queue.BaseJob
    Kime    string `json:"kime"`
    Konu    string `json:"konu"`
    Icerik  string `json:"icerik"`
}

func (j *EpostaGonderJob) Name() string            { return "eposta_gonder" }
func (j *EpostaGonderJob) Queue() string           { return "epostalar" }
func (j *EpostaGonderJob) Retries() int            { return 5 }
func (j *EpostaGonderJob) Timeout() time.Duration  { return 30 * time.Second }

func (j *EpostaGonderJob) Handle(ctx context.Context) error {
    // e-postayı gönder...
    return mailer.Send(ctx, j.Kime, j.Konu, j.Icerik)
}
```

CLI ile job stub'ı üretin:

```bash
zatrano gen job eposta_gonder
# → modules/jobs/eposta_gonder.go
```

### Job Gönderme (Dispatch)

```go
// Uygulama başlangıcında job türlerini kaydedin
app.Queue.Register("eposta_gonder", func() queue.Job { return &jobs.EpostaGonderJob{} })

// Hemen gönder
app.Queue.Dispatch(ctx, &jobs.EpostaGonderJob{
    Kime:   "kullanici@example.com",
    Konu:   "Hoş geldiniz!",
    Icerik: "Merhaba dünya",
})

// Gecikmeli gönder (Redis ZADD sorted set)
app.Queue.Later(ctx, 5*time.Minute, &jobs.EpostaGonderJob{
    Kime: "kullanici@example.com",
    Konu: "Takip",
})
```

### Worker Süreci

Kuyruktan jobları işleyen uzun çalışan bir worker başlatın:

```bash
zatrano queue work
zatrano queue work --queue epostalar --queue bildirimler
zatrano queue work --tries 5 --timeout 120s --sleep 5s
```

Worker otomatik olarak:
- Redis BRPOP ile FIFO sırasında polllar
- Geciktirilmiş jobları her saniye taşır (ZADD → LPUSH)
- Başarısız jobları **üssel geri çekilme** ile yeniden dener (2^deneme saniye)
- Kalıcı olarak başarısız jobları `zatrano_failed_jobs` PostgreSQL tablosuna kaydeder
- `Handle()` içindeki panic’lerden kurtulur
- SIGINT/SIGTERM: önce Redis’ten yeni job alımını keser (BRPOP iptali), eldeki `Handle()` işi bitene kadar bekler (aynı anda tek iş parçacığında en fazla bir job), sonra çıkar; gecikmiş job taşıyan migrate döngüsü için ana `context` süreç sonunda iptal edilir

### Başarısız Joblar

Maksimum yeniden deneme sayısını aşan joblar, hata mesajı, stack trace ve orijinal payload ile PostgreSQL’e kaydedilir.

```bash
# Başarısız jobları listele
zatrano queue failed

# Belirli bir başarısız jobı yeniden dene
zatrano queue retry 42

# Tüm başarısız jobları yeniden dene
zatrano queue retry --all

# Tüm başarısız job kayıtlarını sil
zatrano queue flush
```

**Veritabanı migration:** `zatrano db migrate` çalıştırın — `000003_failed_jobs` migration’ı gerekli tabloyu oluşturur.

### Kuyruk Mimarisi

| Bileşen | Redis Yapısı | Amaç |
|---|---|---|
| Hazır kuyruk | `LIST` (LPUSH/BRPOP) | FIFO job işleme |
| Geciktirilmiş joblar | `SORTED SET` (ZADD) | Zaman bazlı zamanlama |
| Başarısız joblar | Veritabanı tablosu (`zatrano_failed_jobs`, migration) | Kalıcı hata kayıtları |

---

## Mail Sistemi (E-posta)

ZATRANO, HTML şablon desteği, asenkron gönderim için kuyruk entegrasyonu, ek dosya desteği ve yeniden kullanılabilir e-posta tanımları için Mailable deseni sunan **çok sürücülü bir mail sistemi** sağlar.

**Bildirimler:** `core.Bootstrap`, **`App.Notifications`** (`*notifications.Manager`) oluşturur ve **`App.Mail`** üzerinden çalışan **`mail`** kanalını kaydeder. **`pkg/auth`** içindeki **şifre sıfırlama** ve **e-posta doğrulama** gönderimleri **`SendToChannels(..., "mail")`** ile yapılır (`notifications.NewNotification`, isteğe `WithData("html", …)`). Aynı yöneticiye veritabanı, SMS veya push kanalları eklenebilir.

### Yapılandırma

```yaml
# config/dev.yaml
mail:
  driver: smtp          # smtp | log (log = geliştirme/test)
  from_name: "Uygulamam"
  from_email: "noreply@uygulamam.com"
  templates_dir: "views/mails"
  smtp:
    host: smtp.example.com
    port: 587
    username: kullanici
    password: sifre
    encryption: tls     # tls | starttls | ""
```

### E-posta Gönderme

```go
import "github.com/zatrano/zatrano/pkg/mail"

// Basit mesaj
app.Mail.Send(ctx, &mail.Message{
    To:      []mail.Address{{Email: "kullanici@example.com", Name: "Deniz"}},
    Subject: "Hoş Geldiniz!",
    HTMLBody: "<h1>Merhaba Deniz!</h1>",
})

// Şablon ile
app.Mail.SendTemplate(ctx,
    []mail.Address{{Email: "kullanici@example.com"}},
    "Uygulamamıza Hoş Geldiniz",
    "welcome",    // views/mails/welcome.html
    "default",    // views/mails/layouts/default.html
    map[string]any{"Name": "Deniz"},
)

// Kuyruk ile asenkron
app.Mail.Queue(ctx, &mail.Message{
    To:      []mail.Address{{Email: "kullanici@example.com"}},
    Subject: "Bülten",
    HTMLBody: body,
})
```

### Mailable Deseni

Yapılandırılmış, yeniden kullanılabilir e-posta tanımları üretin:

```bash
zatrano gen mail hosgeldiniz
# → modules/mails/hosgeldiniz_mail.go
# → views/mails/hosgeldiniz.html
```

```go
type HosgeldinizMail struct {
    Ad    string
    Email string
}

func (m *HosgeldinizMail) Build(b *mail.MessageBuilder) error {
    b.To(m.Ad, m.Email).
        Subject("Hoş Geldiniz!").
        View("hosgeldiniz", "default", map[string]any{"Name": m.Ad}).
        AttachData("rehber.pdf", pdfBytes, "application/pdf")
    return nil
}

// Senkron gönder
app.Mail.SendMailable(ctx, &mails.HosgeldinizMail{Ad: "Deniz", Email: "deniz@example.com"})

// Kuyruk ile asenkron
app.Mail.QueueMailable(ctx, &mails.HosgeldinizMail{Ad: "Deniz", Email: "deniz@example.com"})
```

### Ek Dosya Desteği

```go
msg := &mail.Message{
    To:      []mail.Address{{Email: "kullanici@example.com"}},
    Subject: "Fatura",
    HTMLBody: body,
    Attachments: []mail.Attachment{
        {Filename: "fatura.pdf", Content: pdfBytes},
        {Filename: "logo.png", Content: logoBytes, Inline: true},
    },
}
app.Mail.Send(ctx, msg)
```

### Şablon Önizleme

Geliştirme sırasında e-posta şablonlarını tarayıcıda önizleyin:

```bash
zatrano mail preview              # şablonları listele
zatrano mail preview welcome      # welcome şablonunu önizle
zatrano mail preview welcome --port 3001
```

Tam yerel mail testi için **Mailpit** veya **MailHog**'u SMTP sunucusu olarak kullanın.

---

## Event / Listener Sistemi

ZATRANO, senkron ve asenkron listener desteği, `ShouldQueue` ile kuyruk tabanlı gönderim ve hızlı geliştirme için üreticiler sunan **merkezi bir event bus**'a (pub/sub) sahiptir.

### Listener Kayıt

```go
// Servis sağlayıcıda / bootstrap'ta (ör. events/event_service_provider.go):

// Senkron (inline fonksiyon)
app.Events.ListenFunc("user.created", func(ctx context.Context, e events.Event) error {
    log.Println("kullanıcı oluşturuldu", e)
    return nil
})

// Struct listener
app.Events.Listen("user.created", &listeners.HosgeldinizMailGonderListener{})

// Birden fazla listener
app.Events.Subscribe("siparis.verildi",
    &listeners.SiparisOnayMailListener{},
    &listeners.StokGuncelleListener{},
)
```

### Event Fırlatma

```go
import "github.com/zatrano/zatrano/pkg/events"

// Event tanımla
type KullaniciOlusturulduEvent struct {
    events.BaseEvent
    KullaniciID uint
    Email       string
}
func (e *KullaniciOlusturulduEvent) Name() string { return "user.created" }

// Senkron fırlat (tüm sync listener'lar tamamlanana kadar bekler)
app.Events.Fire(ctx, &KullaniciOlusturulduEvent{KullaniciID: 1, Email: "ali@example.com"})

// Asenkron fırlat (goroutine'ler, hatalar sadece loglanır)
app.Events.FireAsync(ctx, &KullaniciOlusturulduEvent{KullaniciID: 1, Email: "ali@example.com"})
```

### Kuyruk Tabanlı Listener (`ShouldQueue`)

`ShouldQueue` arayüzünü implemente eden listener'lar kuyruk job'u olarak gönderilir:

```go
type HosgeldinizMailGonderListener struct{}

func (l *HosgeldinizMailGonderListener) Handle(ctx context.Context, event events.Event) error {
    // arka plan worker'ında çalışır
    return nil
}

func (l *HosgeldinizMailGonderListener) Queue() string { return "events" }
func (l *HosgeldinizMailGonderListener) Retries() int  { return 3 }
```

Redis yapılandırılmışsa `ShouldQueue` uygulayan listener'lar otomatik olarak Queue sistemi üzerinden gönderilir.

### Üretici

```bash
zatrano gen event kullanici_olusturuldu
# → modules/events/kullanici_olusturuldu_event.go

zatrano gen listener hosgeldiniz_mail_gonder
# → modules/listeners/hosgeldiniz_mail_gonder_listener.go  (senkron)

zatrano gen listener hosgeldiniz_mail_gonder --queued
# → modules/listeners/hosgeldiniz_mail_gonder_listener.go  (ShouldQueue / asenkron)
```

### Event Servis Sağlayıcı

Tüm listener kayıtlarını tek bir yerde toplayın:

```go
// modules/events/event_service_provider.go
package myevents

import (
    "github.com/zatrano/zatrano/pkg/core"
    "myapp/modules/listeners"
)

// Register tüm event listener'larını bağlar. main veya bootstrap'tan çağırın.
func Register(app *core.App) {
    app.Events.Listen("user.created", &listeners.HosgeldinizMailGonderListener{})
    app.Events.Listen("siparis.verildi", &listeners.SiparisOnayMailListener{})
}
```

---

## Repository / Veri Sistemi

ZATRANO, veri erişimini standartlaştırmak, yeniden kullanılabilir sorgu scope'ları (kapsamlar) dayatmak ve sayfalama ile soft-delete gibi yaygın görevleri otomatikleştirmek için GORM üzerinde **generic repository deseni** sağlar.

### Base Model & Generic Repository

Kayıt kimliği (ID), zaman damgaları (timestamps) ve standart soft-delete davranışı için doğrudan `repository.Model` struct'ını modellerinize gömün.

```go
import "github.com/zatrano/zatrano/pkg/repository"

type User struct {
    repository.Model
    Name  string
    Email string
}

// Servis katmanınızda:
repo := repository.New[User](app.DB)

// Oluştur
repo.Create(ctx, &User{Name: "Alice", Email: "alice@example.com"})

// Soft Delete
repo.DeleteByID(ctx, 1)

// Geri Yükle
repo.Restore(ctx, 1)
```

### Zincirleme Kapsamlar (Chainable Scopes)

GORM iç yapısını handler (işleyici) fonksiyonlarınıza sızdırmadan karmaşık sorgular oluşturun.

```go
// Önceden tanımlanmış kapsamlar (scopes)
scopes := repository.Scopes(
    repository.Active(),
    repository.Where("email LIKE ?", "%@example.com"),
    repository.PreloadAll(), // N+1 sorununu önlemek için hazır yükleme
    repository.OrderBy("created_at DESC"),
    repository.Limit(10),
)

users, _ := repo.FindAll(ctx, scopes...)
```

### Sayfalama (Pagination)

Sayfalama yerleşiktir ve yanıtları standartlaştırır. `repo.Paginate`, öğeleri ve standartlaştırılmış meta verileri içeren bir `Page[T]` döndürür.

```go
opts := repository.PaginateOpts{Page: 1, PerPage: 15}

page, _ := repo.Paginate(ctx, opts, repository.Active())

// page.Items (verileriniz)
// page.Pagination.Total, page.Pagination.CurrentPage, vb.

// HTML şablonlarda kullanmak üzere sayfalama linkleri alma
links := page.Pagination.Links("/users", "&sort=desc")
```

---

### Uluslararasılaştırma (`i18n`)

Uygulama metinleri **`locales_dir`** altında **JSON** dosyalarında tutulur: **`{locales_dir}/{etiket}.json`** (ör. `locales/tr.json`). İç içe nesneler **nokta anahtarlara** düzleştirilir (`app.welcome`).

- **Yapılandırma:** `i18n.enabled`, `i18n.default_locale`, `i18n.supported_locales`, `i18n.locales_dir`, isteğe `i18n.cookie_name` (varsayılan `zatrano_lang`), `i18n.query_key` (varsayılan `lang`). **`i18n.enabled: true`** iken **`locales_dir`** dizini diskte olmalıdır (yüklemede doğrulanır).
- **Çözüm sırası:** sorgu (`?lang=`), çerez, **`Accept-Language`**, ardından **`default_locale`**.
- **Handler:** `github.com/zatrano/zatrano/pkg/i18n` — sabit metin için **`i18n.T(c, "app.welcome")`**; değişkenli çeviri için **`i18n.Tf(c, "app.hello_user", map[string]any{"Name": ad})`** veya struct; JSON içinde **`{{.Name}}`** gibi **`text/template`** ifadeleri. **`map`** ile basit `{{.Alan}}` otomatik uyumludur; Fiber dışında **`Bundle.Format`**. i18n kapalıyken **`T`** / **`Tf`** anahtarı döner (**`Tf`** hata **`nil`**).
- **GET /** yanıtında **`i18n`** nesnesi (`enabled`; açıksa `default_locale`, `supported_locales`, **`active_locale`**).
- **Doğrulama mesajları**, i18n açık olduğunda `validation.*` anahtarlarından otomatik olarak çözülür (bkz. [Validation](#validation-doğrulama)).

---

## View / Template Sistemi

ZATRANO, Go'nun `html/template` paketi üzerine inşa edilmiş, layout kalıtımı, yeniden kullanılabilir bileşen partial'ları, form builder yardımcıları, session tabanlı flash mesajları, eski girdi yeniden doldurma ve Vite/esbuild asset pipeline'ı destekleyen birinci sınıf bir template motoru sunar.

### Layout Kalıtımı (`{{extends}}` / `{{block}}`)

Alt görünümler, **ilk satırda** üst layout'u bildirir:

```html
{{extends "layouts/app"}}

{{block "title"}}Dashboard{{end}}

{{block "content"}}
<h1>Hoşgeldiniz, {{.User.Name}}</h1>
{{end}}
```

Motor, layout'taki `{{block "isim"}}varsayılan{{end}}` tanımlamalarını bulur ve alt görünümün geçersiz kılmalarıyla değiştirir. Alt görünümün tanımlamadığı bloklar layout'taki varsayılan içeriği render eder.

**`layouts/app` içindeki yerleşik bloklar:**

| Blok | Amaç |
|------|------|
| `title` | `<title>` içeriği |
| `head` | Ekstra `<head>` öğeleri (meta, stil) |
| `body_class` | `<body>` üzerindeki CSS sınıfları |
| `header` | Üst navigasyon çubuğu |
| `nav` | Header içindeki navigasyon linkleri |
| `content` | Ana sayfa içeriği (**zorunlu**) |
| `footer` | Sayfa alt bölümü |
| `scripts` | `</body>` öncesi ekstra `<script>` etiketleri |

### Bileşen Sistemi (`views/components/`)

Bileşenler, `views/components/` altındaki sıradan `.html` dosyalarıdır. Başlangıçta otomatik olarak keşfedilir ve adlandırılmış şablonlar olarak kaydedilir — manuel import gerekmez.

```html
{{/* Satır içi uyarı */}}
{{template "components/alert" (dict "Type" "success" "Message" "Kaydedildi!")}}

{{/* Doğrulama hatası ve eski değerle form girdi */}}
{{template "components/form-input" (dict
  "Type"     "email"
  "Name"     "email"
  "Label"    "E-posta Adresi"
  "Value"    (old "email" .Old)
  "Required" true
  "Error"    (index .Errors "email")
)}}
```

**Yerleşik bileşenler:**

| Bileşen | Açıklama |
|---------|---------|
| `components/alert` | SVG ikonlu renkli uyarı kutusu (`success`, `error`, `warning`, `info`) |
| `components/button` | Variant destekli `<button>` (`primary`, `secondary`, `danger`, `ghost`) |
| `components/form-input` | Label + hata/ipucu ile `<input>` |
| `components/form-select` | Seçenek listesi, label, hata ile `<select>` |
| `components/form-textarea` | Label, satır, hata ile `<textarea>` |
| `components/csrf` | Gizli `<input name="_csrf">` |
| `components/pagination` | Offset tabanlı sayfalama linkleri |
| `partials/flash-messages` | Kuyruktaki tüm flash mesajları render eder |

### Form Builder

Ham HTML yazmadan HTML formları oluşturmak için template fonksiyonları:

```html
{{form_open "/kullanicilar" "POST"}}
  {{csrf_field .CSRF}}

  {{input "text" "ad" (old "ad" .Old) `class="form-control"`}}
  {{textarea "bio" (old "bio" .Old) `rows="4"`}}
  {{select "rol" .FormRol (slice (arr "admin" "Yönetici") (arr "user" "Kullanıcı"))}}
  {{checkbox "aktif" "1" true}}

  <button type="submit">Kaydet</button>
{{form_close}}
```

| Yardımcı | İmza | Çıktı |
|----------|------|-------|
| `form_open` | `action method [attrs...]` | `<form ...>` |
| `form_close` | — | `</form>` |
| `csrf_field` | `token` | `<input type="hidden" name="_csrf" value="...">` |
| `input` | `type name value [attrs...]` | `<input ...>` |
| `textarea` | `name value [attrs...]` | `<textarea>...</textarea>` |
| `select` | `name selected [][2]string [attrs...]` | `<select>...</select>` |
| `checkbox` | `name value checked [attrs...]` | `<input type="checkbox" ...>` |

### Flash Mesajları

Yönlendirmeden önce flash mesajı ayarlayın; bir sonraki istekte kullanılabilir ve ardından otomatik olarak temizlenir.

```go
import "github.com/zatrano/zatrano/pkg/view/flash"

// Handler içinde:
flash.Set(c, flash.Success, "Kayıt başarıyla kaydedildi.")
flash.Set(c, flash.Error,   "Bir şeyler yanlış gitti.")
return c.Redirect("/panel")
```

Flash türleri: **`flash.Success`**, **`flash.Error`**, **`flash.Warning`**, **`flash.Info`**.

Layout veya görünümünüzde `{{template "partials/flash-messages" .}}` kuyruktaki tüm mesajları render eder. Flash verisi `a.View.ViewData(c)` tarafından otomatik olarak enjekte edilir.

### Eski Girdi Yardımcısı (`{{old}}`)

Başarısız form doğrulamasından sonra kullanıcının girdisini saklayın, böylece form yeniden doldurulur:

```go
// Handler'da — geri yönlendirmeden önce:
flash.SetOld(c, map[string]string{
    "email": c.FormValue("email"),
    "ad":    c.FormValue("ad"),
})
flash.Set(c, flash.Error, "Doğrulama başarısız.")
return c.Redirect("/kullanicilar/olustur")
```

Template'de:

```html
{{input "email" "email" (old "email" .Old) `placeholder="siz@ornek.com"`}}
```

### Asset Yardımcısı (`{{asset}}`)

`{{asset "yol"}}` bir statik dosya için versiyonlanmış URL döndürür. Çözüm sırası:

1. **Vite/esbuild manifest** — hash'li dosya adını döndürür (ör. `app-a1b2c3.js`)
2. **MD5 dosya hash'i** — `PublicDir`'deki dosyalar için `?v=<hash>` ekler
3. **Düz URL** — `PublicURL/yol`'a geri döner

```html
{{assetLink "css/app.css"}}
{{assetScript "js/app.js"}}
<img src="{{asset "img/logo.png"}}">
```

### Vite / esbuild Entegrasyonu

`view.asset.vite_manifest`'i build çıktı manifest'inize işaret ettirin. `vite build` sonrası `{{asset}}` hash'li dosya adlarını otomatik çözer. Geliştirme sırasında HMR etkinleştirin:

```yaml
# config/dev.yaml
view:
  dev_mode: true
  asset:
    vite_dev_url: http://localhost:5173
    vite_manifest: public/build/.vite/manifest.json
```

```html
{{/* Dev'de Vite client + modül entry point, prod'da hash'li etiketler */}}
{{viteHead "src/main.ts"}}
```

### Handler'larda View Engine Kullanımı

```go
func (h *Handler) Goster(c fiber.Ctx) error {
    data := h.app.View.ViewData(c, fiber.Map{
        "Title": "Panel",
        "User":  kullanici,
    })
    // ViewData otomatik olarak şunları enjekte eder: Flash, Old, OldFn, CSRF, Path, Method
    return c.Render("panel/goster", data)
}
```

### Template Yardımcı Referansı

| Fonksiyon | Örnek | Açıklama |
|-----------|-------|---------|
| `asset` | `{{asset "app.css"}}` | Versiyonlanmış asset URL'i |
| `assetLink` | `{{assetLink "app.css"}}` | `<link rel="stylesheet">` etiketi |
| `assetScript` | `{{assetScript "app.js"}}` | `<script defer>` etiketi |
| `viteHead` | `{{viteHead "src/main.ts"}}` | Vite entry point etiketleri (dev + prod) |
| `old` | `{{old "email" .Old}}` | Önceki form değeri |
| `csrf_field` | `{{csrf_field .CSRF}}` | Gizli CSRF input |
| `dict` | `{{dict "K" "V"}}` | Satır içi `map[string]any` oluştur |
| `safe` | `{{safe .HTML}}` | String'i `template.HTML` olarak işaretle |
| `nl2br` | `{{nl2br .Metin}}` | Yeni satırları `<br>` ile değiştir |
| `default` | `{{default "—" .Deger}}` | Değer sıfır olduğunda yedek |
| `upper` / `lower` / `title` | `{{upper .Ad}}` | String büyük/küçük harf dönüştürücü |
| `json` | `{{json .Veri}}` | Değeri JSON olarak kodla |
| `hasKey` | `{{hasKey .Map "anahtar"}}` | Map'in anahtar içerip içermediğini kontrol et |
| `concat` | `{{concat "a" "b"}}` | String'leri birleştir |
| `iterate` | `{{range iterate 5}}` | 0…n-1 arasında iterasyon |

### Kod Üretimi (`zatrano gen view`)

```bash
# Modül için index.html + show.html scaffold oluştur:
zatrano gen view gonderi

# Ayrıca form scaffold ile create.html + edit.html üret:
zatrano gen view gonderi --with-form

# Özel layout kullan:
zatrano gen view gonderi --with-form --layout layouts/admin

# Kuru çalıştırma — dosya yazmadan yolları yazdır:
zatrano gen view gonderi --with-form --dry-run
```

### Yapılandırma (`view.*`)

```yaml
view:
  root: views              # template kök dizini (varsayılan: views)
  extension: .html         # template dosya uzantısı (varsayılan: .html)
  components_dir: components  # bileşen partial'ları için alt dizin
  layouts_dir: layouts     # layout şablonları için alt dizin
  dev_mode: true           # önbelleği devre dışı bırak; her istekte yenile

  asset:
    public_dir: public          # statik asset'ler için dosya sistemi yolu
    public_url: /public         # asset URL'leri için URL öneki
    vite_manifest: ""           # Vite/esbuild manifest.json yolu
    vite_dev_url: ""            # Vite dev sunucu URL'i (dev_mode'da HMR)
```

---

## Yapılandırma

- **`.env`**, **`config/{env}.yaml`**, **ortam değişkenleri** (ör. `SECURITY_JWT_SECRET`). Çoklu köken veya **`supported_locales`** gibi **listeler** için **YAML** tercih edin.
- Ayrıntı: `migrations_source`, `migrations_dir`, `seeds_dir`, `openapi_path`, **`http.*`**, **`i18n.*`**, `security.*`, `oauth.*` — `config/examples/dev.yaml`.

### Veritabanı migrasyonları (SQL)

- **`migrations_source`:** **`embed`** (varsayılan) — sürüm numaralı `*.up.sql` / `*.down.sql` dosyaları **`pkg/migrations/sql/<sürücü>/`** altında (`postgres`, `mysql`, `sqlite`, `sqlserver`). `zatrano db migrate`, **golang-migrate** + **`embed`/`iofs`** kaynağı ve **`database_driver`** ile aynı sürücüyü kullanır.
- **`file`** — migrasyonları diskteki **`migrations_dir`** dizininden okur (**`zatrano new`** / scaffold projelerinde genelde `migrations_source: file` ve kök `migrations/`).
- **`db migrate`**, **`db rollback`**, **`db tenants …`** komutlarında **`--migrations <dizin>`** o çalıştırma için **dosya** kaynağını zorunlu kılar (embed kullanılmaz).
- **`zatrano gen model`** yalnızca **`pkg/migrations/sql/postgres/`** altına yeni `.up.sql` / `.down.sql` iskeleti yazar; diğer sürücüler için **embed** kullanıyorsanız dosyaları çoğaltıp uyarlamanız gerekir.
- Depo kökündeki **`migrations/`** klasörü isteğe bağlı disk hazırlığı içindir; **`file`** modunda bkz. **`migrations/README.md`**.
- Hata ayıklama: **`zatrano config print`** (tam, maskeli) veya **`zatrano config print --paths-only`** (sohbete yapıştırmaya uygun özet).
- CI: önce **`zatrano config validate -q`** (hızlı YAML/ortam kontrolü), sonra **`zatrano openapi validate --merged`**, veya tam kapı için **`zatrano verify`** (Geliştirme bölümüne bakın).

---

## Geliştirme

```bash
go test ./... -count=1
go fmt ./...
go vet ./...
golangci-lint run
```

**Tek komut kontrol:** `zatrano verify` (veya POSIX **`make verify`**) — `vet`, `test`, birleşik OpenAPI. Yayın öncesi için **`make verify-race`** / **`zatrano verify --race`**. **`make config-validate`**, **`zatrano config validate`** ile aynıdır.

**Canlı yenileme:** [Air](https://github.com/air-verse/air) kurun, `air` çalıştırın (`.air.toml`). Windows'ta çıktı `./tmp/main.exe`.

**Birleşik OpenAPI dosyası:** `make openapi-export` veya `go run ./cmd/zatrano openapi export --output api/openapi.merged.yaml`.

**Ortam kontrolü:** `zatrano doctor` yapılandırma özeti, **`http`** ara katmanı (CORS, rate limit, timeout, gövde boyutu) ve **`config print --paths-only`** ipucu, açıksa **OAuth**, yedekleme için **`pg_dump` / `pg_restore` / `psql`** PATH bilgisi ve bağlantı testlerini gösterir.

**Kod üret:** `gen module` / `gen crud` wire dosyasını günceller ve **`go fmt`** çalıştırır. **`gen wire`** yalnızca patch (ör. **`--skip-wire`** sonrası). **`gen request`** bağımsız olarak form request struct stubs üretir. **`gen policy`** CRUD metotlarıyla (ViewAny, View, Create, Update, Delete, ForceDelete, Restore) `auth.Policy` implementasyonu üretir. Uygulamalarda **`internal/routes/register.go`**, framework deposunda **`pkg/server/register_modules.go`**.

**Sunucu gömme:** `server.Mount(..., server.MountOptions{RegisterRoutes: …})`; `zatrano.StartOptions.RegisterRoutes` üretilen uygulamalarda bu çağrıyı iletir.

---

## Dokümantasyon

- **İngilizce:** [`README.md`](README.md)
- **Türkçe:** bu dosya (`README.tr.md`)

İki dosyayı da aynı değişiklikte güncelleyin.

---

## Katkı

Öneri ve PR'lar memnuniyetle karşılanır. Davranış veya CLI değişikliklerinde **her iki** README'yi de güncelleyin.

---

## Lisans

Belirlenecek.
