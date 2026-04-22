# Quickstart: System Locale Localized Messages

**Feature**: 008-locale-messages

## Build & Run

```bash
# Build the application
go build -o mytets ./cmd/mytets

# Run with default locale (English)
./mytets help

# Run with Ukrainian locale
LANG=uk_UA.UTF-8 ./mytets help

# Run a subcommand — phrases are not localized
LANG=uk_UA.UTF-8 ./mytets one

# Trigger a localized error
LANG=uk_UA.UTF-8 ./mytets --output xml one
```

## Run Tests

```bash
# Unit tests
go test ./internal/i18n/...
go test ./internal/...

# Integration tests
go test ./tests/integration/...

# All tests with race detection
go test -race ./...
```

## Add a New Language

1. Copy the English reference file:
   ```bash
   cp internal/i18n/locales/en.toml internal/i18n/locales/de.toml
   ```

2. Translate the values in `de.toml` (keys stay the same)

3. Rebuild:
   ```bash
   go build -o mytets ./cmd/mytets
   ```

4. Test:
   ```bash
   LANG=de_DE.UTF-8 ./mytets help
   ```

No Go source code changes required.
