# Releasing, semver & changelog

## Conventional Commits

Write commit messages in [Conventional Commits](https://www.conventionalcommits.org/) form so `git-cliff` can group them and compute **semver** bumps:

| Prefix   | Use for |
|----------|---------|
| `feat:`  | new behavior (minor bump in 1.x) |
| `fix:`   | bug fixes (patch) |
| `feat!:` / `BREAKING CHANGE:` in footer | major bump (see Conventional Commits) |
| `docs:`, `chore:`, `ci:`, `refactor:`, `test:`, `perf:`, `build:`, `style:`, `revert:` | as appropriate |

**Scopes** are optional: `feat(api): add handler`.

**Pull request titles** are validated in CI (Conventional-Commits style). The merge commit (or squashed) message on `main` should follow the same rules.

## `CHANGELOG.md`

Regenerate from git history and `cliff.toml`:

```bash
make changelog
```

Requires [`git-cliff`](https://github.com/orhun/git-cliff#installation) on your `PATH`.

## Next version (without tagging)

```bash
make next-version
# same as: git-cliff --bumped-version
```

This uses conventional commits **since the last `v*.*.*` tag** and `cliff.toml` `[bump]` rules. Empty output means nothing to release or no bumpable commits.

## Tagging a release (semver)

1. Ensure `main` is green and commits follow conventions.
2. (Optional) `make changelog` and commit `CHANGELOG.md`, or use the [Changelog PR workflow](https://github.com/zatrano/zatrano/actions) (`changelog-pr.yml`, **Run workflow**).
3. Create an annotated tag:

   ```bash
   git tag -a v0.0.2 -m "v0.0.2: short description"
   git push origin v0.0.2
   ```

4. The **Release** workflow publishes a [GitHub Release](https://github.com/zatrano/zatrano/releases) with notes from `git-cliff` for that tag.

## Automation summary

| Workflow | Trigger | Effect |
|----------|---------|--------|
| `conventional-commits-pr.yml` | pull requests | Checks PR title against Conventional Commits |
| `changelog-pr.yml` | `workflow_dispatch` | Regenerates `CHANGELOG.md` and opens a PR |
| `release.yml` | push `v*` tag | Publishes GitHub Release with `git-cliff` notes for the tag |

## `go install` and tags

Most Zatrano users run:

- **`go install ...@latest`** — highest **published** `v*.*.*` tag the module proxy has (what people expect for “latest”).
- **`go install ...@vX.Y.Z`** — exact, reproducible version.
- **`go install ...@main`:** tip of the default branch; use if you need a commit **before** the next tag, or for local development.

Tag releases (see [Automation summary](#automation-summary)) so `@latest` always tracks real releases. If the proxy lags, use `GOPROXY=direct` or `go clean -modcache` before reinstalling.

See also [README.md](../README.md#installation).
