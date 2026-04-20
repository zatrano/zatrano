# Migrations (disk)

ZATRANO’s built-in schema lives in **`pkg/migrations/sql/<driver>/`** and is applied by default (`migrations_source: **embed**` in config).

Use this directory only when **`migrations_source: file`** (or `zatrano db migrate --migrations ./migrations`). New apps from `zatrano scaffold` set `migrations_source: file` and keep starter SQL here.
