# CLI Locale Contract: System Locale Localized Messages

**Feature**: 008-locale-messages  
**Date**: 2026-04-21

## Overview

This contract defines the expected CLI behavior when the system locale affects output language. The application detects the system locale and renders all user-facing text (help, errors, structural labels) in the corresponding language. Phrase content from `phrases.json` is never affected.

## Locale Detection

**Input**: Environment variables (read in priority order)

| Priority | Variable | Example |
|----------|----------|---------|
| 1 | `LC_ALL` | `uk_UA.UTF-8` |
| 2 | `LC_MESSAGES` | `uk_UA.UTF-8` |
| 3 | `LANG` | `uk_UA.UTF-8` |

**Language code extraction**: First segment before `_` or `.` (e.g., `uk_UA.UTF-8` → `uk`)

**Fallback**: If no variable is set, or the extracted language has no matching `.toml` file, the application uses English (`en`).

## Localized Help Output

### Root Command Help

**Command**: `mytets help` or `mytets --help`

**English locale** (`LANG=en_US.UTF-8`):
```
mytets — a lightweight CLI tool

Usage:
  mytets [command]

Available Commands:
  list        Display a list of random phrases
  one         Display one random phrase
  version     Print the application version and exit

Flags:
  -h, --help   help for mytets

Use "mytets [command] --help" for more information about a command.
```

**Ukrainian locale** (`LANG=uk_UA.UTF-8`):
```
mytets — легкий інструмент командного рядка

Використання:
  mytets [команда]

Доступні команди:
  list        Показати список випадкових фраз
  one         Показати одну випадкову фразу
  version     Вивести версію програми та завершити

Прапори:
  -h, --help   довідка для mytets

Використовуйте "mytets [команда] --help" для отримання додаткової інформації про команду.
```

### Subcommand Help (example: `mytets one --help`)

**Ukrainian locale** (`LANG=uk_UA.UTF-8`):
```
Команда one виводить випадкову фразу у текстовому або JSON форматі.

Використання:
  mytets one [прапори]

Прапори:
  -h, --help            довідка для one
  -o, --output string   Формат виводу: "text" (за замовчуванням) або "json" (за замовчуванням "text")
```

## Localized Error Messages

### Error: No Phrases Available

**Command**: `mytets one` (when phrase source is empty)

| Locale | Error message (stderr) |
|--------|----------------------|
| `en` | `no phrases available` |
| `uk` | `фрази відсутні` |

### Error: Invalid Output Format

**Command**: `mytets --output xml one`

| Locale | Error message (stderr) |
|--------|----------------------|
| `en` | `unsupported output format: "xml"` |
| `uk` | `непідтримуваний формат виводу: "xml"` |

### Error: Unknown Command

**Command**: `mytets foo`

| Locale | Error message (stderr) |
|--------|----------------------|
| `en` | `unknown command "foo" for "mytets"` |
| `uk` | `невідома команда "foo" для "mytets"` |

## Non-Localized Content

The following outputs are **never** affected by locale:

| Output | Reason |
|--------|--------|
| Phrase text from `phrases.json` | Authored content, not UI messages (FR-008) |
| Version string (e.g., `1.0.0`) | Technical identifier, not natural language |
| JSON field names (`"message"`, `"version"`) | Machine-readable contract, locale-independent |
| Exit codes | Numeric, locale-independent |

## Localization File Format

Each language has one TOML file: `internal/i18n/locales/{lang}.toml`

**Example** (`en.toml`):
```toml
# Root command
root.short = "mytets — a lightweight CLI tool"

# Subcommands
version.short = "Print the application version and exit"
one.short = "Display one random phrase"
one.long = "The one command outputs a random phrase in plain text or JSON format."
list.short = "Display a list of random phrases"
list.long = "The list command outputs multiple unique random phrases in plain text or JSON format."

# Flags
flag.output = "Output format: \"text\" (default) or \"json\""
flag.count = "Number of phrases to return"

# Cobra structural labels
help.usage = "Usage:"
help.available_commands = "Available Commands:"
help.flags = "Flags:"
help.global_flags = "Global Flags:"
help.additional_help = "Use \"%s [command] --help\" for more information about a command."
help.aliases = "Aliases:"
help.help_for = "help for %s"

# Errors
error.no_phrases = "no phrases available"
error.invalid_output_format = "unsupported output format: \"%s\""
error.failed_select_phrase = "failed to select phrase: %w"
error.failed_encode_json = "failed to encode JSON: %w"
error.unknown_command = "unknown command \"%s\" for \"%s\""
```

## Adding a New Language

1. Copy `en.toml` to `{lang}.toml` (e.g., `de.toml` for German)
2. Translate all values (keys remain unchanged)
3. Run `go build` — the new file is auto-embedded via `//go:embed locales/*.toml`
4. No Go source code changes required
