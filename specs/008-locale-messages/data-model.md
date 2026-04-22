# Data Model: System Locale Localized Messages

**Feature**: 008-locale-messages  
**Date**: 2026-04-21

## Entities

### Bundle

The central localization data structure holding all loaded translations.

| Field | Type | Description |
|-------|------|-------------|
| translations | map[string]map[string]string | Two-level map: language code → translation key → translated string |
| lang | string | The active language code resolved from the system locale (e.g., `uk`, `en`) |

**Invariants**:
- The `en` key is always present in `translations` (loaded from `en.toml`)
- `lang` is always a valid key in `translations`, or falls back to `en`
- After initialization, `translations` is read-only (no concurrent mutation)

### Localization File (TOML)

A flat key-value TOML file containing all translatable strings for one language.

| Field | Type | Description |
|-------|------|-------------|
| (key) | string | Dot-separated identifier (e.g., `root.short`, `error.no_phrases`) |
| (value) | string | The translated string for the given key in this language |

**Naming convention**: `{language_code}.toml` (e.g., `en.toml`, `uk.toml`)

**Key naming convention** (dot-separated hierarchy):
- `root.short` — root command short description
- `root.long` — root command long description (if any)
- `version.short` — version command short description
- `one.short` — one command short description
- `one.long` — one command long description
- `list.short` — list command short description
- `list.long` — list command long description
- `flag.output` — output flag description
- `flag.count` — count flag description
- `error.no_phrases` — "no phrases available" error
- `error.invalid_output_format` — invalid output format error
- `error.failed_select_phrase` — "failed to select phrase" error
- `error.failed_encode_json` — "failed to encode JSON" error
- `help.usage` — "Usage:" label
- `help.available_commands` — "Available Commands:" label
- `help.flags` — "Flags:" label
- `help.global_flags` — "Global Flags:" label
- `help.additional_help` — "Use ... for more information" text
- `help.aliases` — "Aliases:" label
- `error.unknown_command` — "unknown command" Cobra error
- `error.unknown_flag` — "unknown flag" Cobra error
- `root.long` — root command long description (reserved for future use)

### Locale

The resolved system language setting.

| Field | Type | Description |
|-------|------|-------------|
| raw | string | The raw value from the environment variable (e.g., `uk_UA.UTF-8`) |
| lang | string | The extracted language code (e.g., `uk`) |

**Resolution rules** (POSIX priority):
1. Read `LC_ALL` — if non-empty, use it
2. Read `LC_MESSAGES` — if non-empty, use it
3. Read `LANG` — if non-empty, use it
4. Default to `en`

**Parsing**: Extract language code by splitting on `_` or `.`, taking the first segment. Empty or `C`/`POSIX` values map to `en`.

## Relationships

```
Bundle 1──* Localization File    (one bundle loads many TOML files)
Bundle 1──1 Locale               (one active locale per execution)
Localization File 1──* Translatable String  (each file has all keys)
```

## State Transitions

The localization system has no state transitions after initialization. The flow is:

1. **Startup**: Detect locale → resolve language code → load all TOML bundles from embedded FS → set active language
2. **Runtime**: `Translate(key)` performs read-only map lookups (active lang → `en` fallback → key name)

No mutation occurs after step 1.
