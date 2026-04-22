package i18n

import (
	"embed"
	"io/fs"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

//go:embed locales/*.toml
var localesFS embed.FS

var (
	bundle map[string]map[string]string
	lang   string
)

// LoadBundle reads all TOML files from the embedded locales/ directory and
// populates the bundle map. Language codes are derived from filenames
// (e.g. "uk.toml" → "uk"). Must be called once at startup before Translate.
func LoadBundle() {
	bundle = make(map[string]map[string]string)

	entries, err := fs.ReadDir(localesFS, "locales")
	if err != nil {
		return
	}

	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".toml") {
			continue
		}
		code := strings.TrimSuffix(name, ".toml")

		data, err := fs.ReadFile(localesFS, "locales/"+name)
		if err != nil {
			continue
		}

		m := make(map[string]string)
		if err := toml.Unmarshal(data, &m); err != nil {
			continue
		}
		bundle[code] = m
	}
}

// DetectLocale reads LC_ALL → LC_MESSAGES → LANG from the environment,
// extracts the language code, and sets it as the active language. If the
// detected language has no matching bundle, falls back to "en".
func DetectLocale() {
	raw := os.Getenv("LC_ALL")
	if raw == "" {
		raw = os.Getenv("LC_MESSAGES")
	}
	if raw == "" {
		raw = os.Getenv("LANG")
	}

	code := extractLangCode(raw)
	if _, ok := bundle[code]; ok {
		lang = code
	} else {
		lang = "en"
	}
}

// Translate returns the translated string for key in the active language.
// Fallback order: active language → "en" → key name.
func Translate(key string) string {
	if m, ok := bundle[lang]; ok {
		if v, ok := m[key]; ok {
			return v
		}
	}
	if m, ok := bundle["en"]; ok {
		if v, ok := m[key]; ok {
			return v
		}
	}
	return key
}

// SetLang overrides the active language. Intended for testing.
func SetLang(l string) {
	lang = l
}

// Lang returns the current active language code.
func Lang() string {
	return lang
}

// BundleLen returns the number of loaded language bundles.
func BundleLen() int {
	return len(bundle)
}

// extractLangCode splits a locale string like "uk_UA.UTF-8" into the
// language code "uk". Empty, "C", and "POSIX" values return "en".
func extractLangCode(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "C" || raw == "POSIX" {
		return "en"
	}
	// Split on '_' or '.' and take the first segment.
	for _, sep := range []string{"_", "."} {
		if i := strings.Index(raw, sep); i > 0 {
			return raw[:i]
		}
	}
	return raw
}
