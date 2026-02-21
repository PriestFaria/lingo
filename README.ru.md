# lingo

Статический анализатор лог-сообщений для Go.  
Находит нарушения стиля и потенциальные утечки данных в вызовах `log`, `log/slog` и `go.uber.org/zap`.

[![CI](https://github.com/PriestFaria/lingo/actions/workflows/ci.yml/badge.svg)](https://github.com/PriestFaria/lingo/actions/workflows/ci.yml)
&nbsp;[🇬🇧 English](README.md)

---

## Правила

| #   | Правило                                                              | Пример нарушения               |
| --- | -------------------------------------------------------------------- | ------------------------------ |
| 1   | Сообщение должно начинаться со **строчной** буквы                    | `log.Info("Starting server")`  |
| 2   | Сообщение должно быть на **английском** языке                        | `log.Info("запуск сервера")`   |
| 3   | Нет **эмодзи** и повторяющейся пунктуации (`!!`, `...`)              | `log.Info("done! 🚀")`         |
| 4   | Нет ключевых слов **чувствительных данных** (`password`, `token`, …) | `log.Info("user token: " + t)` |

Правило 1 (строчная буква) поддерживает **авто-исправление** через `suggested fixes`.

## Поддерживаемые логгеры

- `log` (стандартная библиотека)
- `log/slog` (стандартная библиотека, Go 1.21+)
- `go.uber.org/zap`

Форматные методы (`Printf`, `Infof`, …) поддерживаются полностью.

---

## Установка

**Требования:** Go 1.22+

```bash
go install github.com/PriestFaria/lingo/cmd/lingo@latest
```

Бинарник будет доступен как `lingo` в `$(go env GOPATH)/bin`.

---

## Использование

### Standalone через `go vet`

```bash
# Все фильтры включены, настройки по умолчанию
go vet -vettool=$(go env GOPATH)/bin/lingo ./...

# С конфигурационным файлом
go vet -vettool=$(go env GOPATH)/bin/lingo -config=.lingo.json ./...
```

### Плагин для golangci-lint (Linux / macOS)

**1. Клонировать и собрать плагин**

```bash
git clone https://github.com/PriestFaria/lingo.git
cd lingo
go build -buildmode=plugin -o /path/to/your/project/lingo.so ./plugin/
```

> Плагин требует сборки из исходников — это ограничение Go plugin system.

**2. Конфигурация `.golangci.yml`**

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

**3. Запуск**

```bash
golangci-lint run ./...
```

> **Примечание:** Go plugin system (`-buildmode=plugin`) не поддерживает Windows.

---

## Конфигурация

lingo настраивается через файл `.lingo.json` или inline в `.golangci.yml`.

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

Все поля опциональны. Отсутствующий фильтр считается **включённым**.  
Чтобы отключить фильтр, задайте явно `false`.

### Приоритет конфигурации (плагин golangci-lint)

1. **Inline** — ключи `filters` / `security` внутри `settings:` в `.golangci.yml`
2. **Файл** — `settings.config: path/to/.lingo.json`
3. **Default** — все фильтры включены, без дополнительных ключевых слов

### Встроенные ключевые слова

`password`, `passwd`, `pass`, `secret`, `token`, `apikey`, `api_key`, `auth`, `credential`, `cred`, `private`, `privkey`, `jwt`, `key`

Кастомные слова добавляются поверх через `extra_keywords`.

---

## Примеры

```go
// ❌ нарушения, которые найдёт lingo

log.Info("Starting server on port 8080")    // должна быть строчная буква
log.Info("запуск сервера")                  // только английский
log.Info("server started 🚀")               // нет эмодзи
log.Error("connection failed!!!")           // нет повторяющейся пунктуации
log.Info("user password: " + password)      // чувствительные данные в литерале
log.Debug("api key", zap.String("key", k))  // чувствительное имя переменной
```

```go
// ✅ корректное использование

log.Info("starting server on port 8080")
log.Info("starting server")
log.Info("server started")
log.Error("connection failed")
log.Info("user authenticated successfully")
log.Debug("api request completed")
```

---

## Тестирование

```bash
# Unit-тесты + analysistest
go test ./internal/...

# End-to-end тесты (собирает бинарник, запускает go vet против 6 sample-проектов)
go test -tags e2e ./test/e2e/
```

### Покрытие

| Пакет               | Покрытие |
| ------------------- | -------- |
| `internal/filters`  | 100%     |
| `internal/analyzer` | 94.1%    |
| `internal/config`   | 92.3%    |

---

## Структура проекта

```
cmd/lingo/             — standalone-бинарник (go vet -vettool)
plugin/                — Go-плагин для golangci-lint
internal/
  analyzer/            — обход AST, роутинг на хэндлеры
  filters/             — реализации правил (FirstLetter, English, Emoji, Security)
  config/              — загрузка .lingo.json и настройки по умолчанию
test/e2e/              — end-to-end тесты против sample-проектов
```

---

## CI

GitHub Actions запускает два джоба при каждом push и pull request в `main`:

- **Unit & Integration** — `go test ./internal/...`
- **E2E** — `go test -tags e2e ./test/e2e/`
