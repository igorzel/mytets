# Feature Specification: System Locale Localized Messages

**Feature Branch**: `008-locale-messages`  
**Created**: 2026-04-21  
**Status**: Draft  
**Input**: User description: "The application should respect system locale and print localized messages/errors/help. When the current system locale is uk_UA.UTF-8 then the application should print information/help/error messages in Ukrainian."

## Clarifications

### Session 2026-04-21

- Q: Should Cobra's own structural labels (`Usage:`, `Flags:`, `Available Commands:`, etc.) and built-in error messages (`unknown command`) also be localized, or only application-defined descriptions? → A: Localize all Cobra structural labels and built-in errors to Ukrainian as well, so the entire help output is consistently in one language.
- Q: What format should the per-language localization files use? → A: TOML — one `.toml` file per language with flat key-value string pairs.
- Q: Should English strings also be extracted into an `en.toml` localization file, or remain hardcoded in the source as the default fallback? → A: Extract English into `en.toml` — all languages including English are managed via TOML files. The `en.toml` file also serves as the reference template for translators.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Help Output in Ukrainian When System Locale Is Ukrainian (Priority: P1)

A user whose system locale is set to `uk_UA.UTF-8` runs `mytets help` or `mytets --help` and expects all help text — command descriptions, flag descriptions, and usage information — to be displayed in Ukrainian.

**Why this priority**: Help text is the most visible user-facing output and directly determines whether a Ukrainian-speaking user can understand and navigate the application. This is the core deliverable.

**Independent Test**: Can be fully tested by setting the system locale to `uk_UA.UTF-8`, running `mytets help`, `mytets --help`, `mytets one --help`, `mytets list --help`, and `mytets version --help`, then verifying that all descriptions and usage text are printed in Ukrainian.

**Acceptance Scenarios**:

1. **Given** the system locale is `uk_UA.UTF-8`, **When** the user runs `mytets help`, **Then** the root command description and all subcommand summaries are displayed in Ukrainian
2. **Given** the system locale is `uk_UA.UTF-8`, **When** the user runs `mytets one --help`, **Then** the short description, long description, and flag descriptions for the `one` command are displayed in Ukrainian
3. **Given** the system locale is `uk_UA.UTF-8`, **When** the user runs `mytets list --help`, **Then** the short description, long description, flag descriptions (including `--count` and `--output`), and usage text are displayed in Ukrainian
4. **Given** the system locale is `uk_UA.UTF-8`, **When** the user runs `mytets version --help`, **Then** the version command description and flag descriptions are displayed in Ukrainian

---

### User Story 2 - Error Messages in Ukrainian When System Locale Is Ukrainian (Priority: P1)

A user whose system locale is set to `uk_UA.UTF-8` triggers an error condition and expects the error message to be displayed in Ukrainian.

**Why this priority**: Error messages are critical user-facing text. A user who cannot understand an error cannot recover from it. This is equally important as help text.

**Independent Test**: Can be fully tested by setting the system locale to `uk_UA.UTF-8`, triggering known error conditions (e.g., providing an invalid output format, or encountering "no phrases available"), and verifying the error messages are in Ukrainian.

**Acceptance Scenarios**:

1. **Given** the system locale is `uk_UA.UTF-8`, **When** the user triggers a "no phrases available" error, **Then** the error message is displayed in Ukrainian
2. **Given** the system locale is `uk_UA.UTF-8`, **When** the user provides an invalid output format, **Then** the error message about the unsupported format is displayed in Ukrainian
3. **Given** the system locale is `uk_UA.UTF-8`, **When** any application-defined error occurs, **Then** the error text is in Ukrainian

---

### User Story 3 - Default English Output When Locale Is Not Ukrainian (Priority: P2)

A user whose system locale is set to any locale other than Ukrainian (e.g., `en_US.UTF-8`, `de_DE.UTF-8`, or an unset locale) runs the application and expects all messages to appear in English, which is the current default language.

**Why this priority**: The application must gracefully fall back to English for any unsupported locale. This protects the existing user experience and ensures the localization mechanism does not break non-Ukrainian environments.

**Independent Test**: Can be fully tested by setting the system locale to `en_US.UTF-8` (or unsetting it), running `mytets help` and triggering errors, and verifying all output is in English identical to the current behavior.

**Acceptance Scenarios**:

1. **Given** the system locale is `en_US.UTF-8`, **When** the user runs `mytets help`, **Then** help text is displayed in English
2. **Given** the system locale is not set or is an unsupported locale, **When** the user runs any command, **Then** all messages are in English
3. **Given** the system locale is `de_DE.UTF-8` (unsupported), **When** the user runs any command, **Then** all messages fall back to English

---

### User Story 4 - Phrase Content Is Not Localized (Priority: P1)

A user runs `mytets one` or `mytets list` and expects the phrase content from `phrases.json` to be displayed exactly as stored, regardless of the system locale.

**Why this priority**: This is a critical constraint. The localization mechanism must not alter the content of phrases. Phrases are authored content, not UI messages.

**Independent Test**: Can be fully tested by running `mytets one` and `mytets list` with the system locale set to `uk_UA.UTF-8` and verifying that the phrase text matches the original `phrases.json` content verbatim.

**Acceptance Scenarios**:

1. **Given** the system locale is `uk_UA.UTF-8`, **When** the user runs `mytets one`, **Then** the phrase text is from `phrases.json` and is not translated
2. **Given** the system locale is `uk_UA.UTF-8`, **When** the user runs `mytets list`, **Then** all phrase texts are from `phrases.json` and are not translated

---

### User Story 5 - Adding a New Language Without Code Changes (Priority: P2)

A contributor wants to add support for a new language (e.g., German). They should be able to do so by creating a single localization file for that language, without modifying any application source code, and then rebuilding the application.

**Why this priority**: Extensibility is a stated design goal. The localization architecture must support adding languages through localization files alone, enabling non-developers to contribute translations.

**Independent Test**: Can be verified by creating a new localization file following the documented format, rebuilding the application, setting the system locale to the new language, and confirming that messages appear in that language.

**Acceptance Scenarios**:

1. **Given** a contributor creates a localization file for a new language following the established format, **When** the application is rebuilt, **Then** messages in that language are available when the matching locale is active
2. **Given** a new language localization file is added, **When** the contributor inspects the change set, **Then** no application source code files were modified — only the new localization file was added
3. **Given** translations need to be outsourced, **When** a translator receives the localization file for a specific language, **Then** they can complete the translations without access to the application source code because all translatable strings for that language are in one self-contained file

---

### Edge Cases

- What happens when the system locale environment variable is malformed or partially set (e.g., `LANG=uk` without the full `uk_UA.UTF-8`)?  
  The application should attempt a best-effort match (matching the language portion `uk`) and fall back to English if no match is found.
- What happens when the system locale is set to a variant of Ukrainian not exactly matching `uk_UA.UTF-8` (e.g., `uk_UA.ISO-8859-5` or just `uk`)?  
  The application should match on the language code (`uk`) regardless of encoding or country variant.
- What happens when locale detection fails entirely (e.g., all locale environment variables are empty)?  
  The application defaults to English.
- What happens if a localization file is incomplete (some keys are missing translations)?  
  The application falls back to the English default for any missing translation key.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The application MUST detect the current system locale at startup by reading standard locale environment variables (such as `LANG`, `LC_ALL`, `LC_MESSAGES`)
- **FR-002**: When the detected locale matches Ukrainian (`uk`), the application MUST display all help text (command descriptions, flag descriptions, usage information) in Ukrainian, including Cobra's structural labels (`Usage:`, `Flags:`, `Available Commands:`, etc.)
- **FR-003**: When the detected locale matches Ukrainian (`uk`), the application MUST display all application-defined error messages in Ukrainian, including Cobra's built-in error messages (e.g., `unknown command`)
- **FR-004**: When the detected locale does not match any supported language, the application MUST fall back to the English TOML file (`en.toml`) for all messages
- **FR-005**: Localization files MUST be compiled into the application binary (embedded at build time) so that no external files are needed at runtime
- **FR-006**: Each supported language MUST have its translations in a single, self-contained TOML file (`.toml`) with flat key-value string pairs, so that the file can be handed to a translator without requiring access to the source code
- **FR-007**: Adding support for a new language MUST NOT require changes to application source code — only adding a new localization file and rebuilding
- **FR-008**: The content of `phrases.json` MUST NOT be affected by localization. Phrases are displayed exactly as stored regardless of the active locale
- **FR-009**: Locale matching MUST be based on the language code (e.g., `uk`) and not require an exact match on the full locale string (e.g., `uk_UA.UTF-8`)
- **FR-010**: The application MUST fall back to the English value from `en.toml` for any individual translation key that is missing from the active language's localization file
- **FR-011**: English strings MUST be managed via an `en.toml` localization file (not hardcoded in source code), which also serves as the reference template for translators creating new language files

### Key Entities

- **Locale**: The system language setting detected from environment variables, resolved to a language code (e.g., `uk`, `en`)
- **Localization File**: A single TOML file (`.toml`) per language containing all translatable UI strings as flat key-value pairs (help text, error messages, flag descriptions). Compiled into the binary
- **Translatable String**: Any user-facing message produced by the application itself (help text, error messages, usage information) — excludes phrase content from `phrases.json`

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: When the system locale is `uk_UA.UTF-8`, 100% of application help text is displayed in Ukrainian
- **SC-002**: When the system locale is `uk_UA.UTF-8`, 100% of application-defined error messages are displayed in Ukrainian
- **SC-003**: When the system locale is any unsupported language, 100% of messages fall back to English with no visible degradation
- **SC-004**: Phrase content from `phrases.json` remains unchanged regardless of the active locale — zero phrase texts are altered
- **SC-005**: A new language can be added by creating exactly one localization file and rebuilding, with zero application source code modifications
- **SC-006**: A translator can complete a language translation using only the localization file, without access to any application source code
- **SC-007**: All existing unit and integration tests continue to pass without modification (backward compatibility)

## Assumptions

- The system locale is determined via standard environment variables (`LC_ALL`, `LC_MESSAGES`, `LANG`) following POSIX conventions; the application does not need to query the OS via system calls
- English is the default and fallback language; English strings are extracted into `en.toml` and serve as the baseline reference for all translations
- The localization mechanism will follow established best practices for Go CLI application internationalization
- Only Ukrainian (`uk`) is required for the initial release; the architecture supports additional languages
- The application runs on Linux, macOS, and Windows; locale detection covers all three platforms via environment variables
- Localization applies to all user-visible text, including Cobra's structural labels (`Usage:`, `Flags:`, `Available Commands:`) and Cobra's built-in error messages (`unknown command`, etc.), to ensure a consistent single-language experience
