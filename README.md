# lingo

Static analyzer for Go log messages.  
Catches style violations and potential data leaks in calls to `log`, `log/slog`, and `go.uber.org/zap`.

[![CI](https://github.com/PriestFaria/lingo/actions/workflows/ci.yml/badge.svg)](https://github.com/PriestFaria/lingo/actions/workflows/ci.yml)
&nbsp;[🇷🇺 Русский](README.ru.md)

---

## Rules

| #   | Rule                                                           | Example violation              |
| --- | -------------------------------------------------------------- | ------------------------------ |
| 1   | Message must start with a **lowercase** letter                 | `log.Info("Starting server")`  |
| 2   | Message must be in **English**                                 | `log.Info("запуск сервера")`   |
| 3   | No **emoji** or repeated punctuation (`!!`, `...`)             | `log.Info("done! 🚀")`         |
| 4   | No **sensitive data** keywords (`password`, `token`, `key`, …) | `log.Info("user token: " + t)` |

Rule 1 (first letter) supports **auto-fix** via `suggested fixes`.

## Supported loggers

- `log` (standard library)
- `log/slog` (standard library, Go 1.21+)
- `go.uber.org/zap`

Format methods (`Printf`, `Infof`, …) are fully supported.

---

## Installation

**Requirements:** Go 1.22+

```bash
go install github.com/PriestFaria/lingo/cmd/lingo@latest
```

The binary will be available as `lingo` in `$(go env GOPATH)/bin`.

---

## Usage

### Standalone via `go vet`

```bash
# All filters enabled, default settings
go vet -vettool=$(go env GOPATH)/bin/lingo ./...

# With a config file
go vet -vettool=$(go env GOPATH)/bin/lingo -config=.lingo.json ./...
```

### golangci-lint plugin (Linux / macOS)

**1. Clone and build the plugin**

```bash
git clone https://github.com/PriestFaria/lingo.git
cd lingo
go build -buildmode=plugin -o /path/to/your/project/lingo.so ./plugin/
```

> The plugin requires building from source — this is a limitation of the Go plugin system.

**2. Configure `.golangci.yml`**

```yaml
version: "2"

linters:
  enable:
    - lingo

linters-settings:
  custom:
    lingo:
      type: goplugin
      path: ./lingo.so
      settings:
        filters:
          first_letter: true
          english: true
          emoji: true
          security: true
        security:
          extra_keywords:
            - cvv
            - ssn
            - otp
```

**3. Run**

```bash
golangci-lint run ./...
```

> **Note:** Go plugin system (`-buildmode=plugin`) does not support Windows.

---

## Configuration

lingo is configured via a `.lingo.json` file or inline in `.golangci.yml`.

### `.lingo.json`

```json
{
  "filters": {
    "first_letter": true,
    "english": true,
    "emoji": true,
    "security": true
  },
  "security": {
    "extra_keywords": ["cvv", "ssn", "otp"]
  }
}
```

All fields are optional. An absent filter defaults to **enabled**.  
Set a filter to `false` to disable it explicitly.

### Config resolution priority (golangci-lint plugin)

1. **Inline** — `filters` / `security` keys inside `settings:` in `.golangci.yml`
2. **File** — `settings.config: path/to/.lingo.json`
3. **Default** — all filters on, no extra keywords

### Built-in sensitive keywords

`password`, `passwd`, `pass`, `secret`, `token`, `apikey`, `api_key`, `auth`, `credential`, `cred`, `private`, `privkey`, `jwt`, `key`

Custom keywords are added on top via `extra_keywords`.

---

## Examples

```go
// ❌ violations lingo will report

log.Info("Starting server on port 8080")    // must start with lowercase
log.Info("запуск сервера")                  // must be in English
log.Info("server started 🚀")               // no emoji
log.Error("connection failed!!!")           // no repeated punctuation
log.Info("user password: " + password)      // sensitive data in literal
log.Debug("api key", zap.String("key", k))  // sensitive variable name
```

```go
// ✅ correct usage

log.Info("starting server on port 8080")
log.Info("starting server")
log.Info("server started")
log.Error("connection failed")
log.Info("user authenticated successfully")
log.Debug("api request completed")
```

---

## Testing

```bash
# Unit tests + analysistest
go test ./internal/...

# End-to-end tests (builds the binary, runs go vet against 6 sample projects)
go test -tags e2e ./test/e2e/
```

### Coverage

| Package             | Coverage |
| ------------------- | -------- |
| `internal/filters`  | 100%     |
| `internal/analyzer` | 94.1%    |
| `internal/config`   | 92.3%    |

---

## Project structure

```
cmd/lingo/             — standalone binary (go vet -vettool)
plugin/                — golangci-lint Go plugin
internal/
  analyzer/            — AST traversal, routing to handlers
  filters/             — rule implementations (FirstLetter, English, Emoji, Security)
  config/              — .lingo.json loading and defaults
test/e2e/              — end-to-end tests against sample projects
```

---

## CI

GitHub Actions runs two jobs on every push and pull request to `main`:

- **Unit & Integration** — `go test ./internal/...`
- **E2E** — `go test -tags e2e ./test/e2e/`
