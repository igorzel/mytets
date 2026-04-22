package i18n

import (
	"testing"
)

func TestLoadBundleDiscoversBothLocales(t *testing.T) {
	LoadBundle()
	if BundleLen() < 2 {
		t.Fatalf("expected at least 2 bundles, got %d", BundleLen())
	}
	for _, code := range []string{"en", "uk"} {
		if _, ok := bundle[code]; !ok {
			t.Errorf("bundle missing language %q", code)
		}
	}
}

func TestDetectLocaleLC_ALL(t *testing.T) {
	LoadBundle()
	t.Setenv("LC_ALL", "uk_UA.UTF-8")
	t.Setenv("LC_MESSAGES", "en_US.UTF-8")
	t.Setenv("LANG", "en_US.UTF-8")
	DetectLocale()
	if Lang() != "uk" {
		t.Errorf("expected lang 'uk', got %q", Lang())
	}
}

func TestDetectLocaleLC_MESSAGES(t *testing.T) {
	LoadBundle()
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "uk_UA.UTF-8")
	t.Setenv("LANG", "en_US.UTF-8")
	DetectLocale()
	if Lang() != "uk" {
		t.Errorf("expected lang 'uk', got %q", Lang())
	}
}

func TestDetectLocaleLANG(t *testing.T) {
	LoadBundle()
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "")
	t.Setenv("LANG", "uk_UA.UTF-8")
	DetectLocale()
	if Lang() != "uk" {
		t.Errorf("expected lang 'uk', got %q", Lang())
	}
}

func TestDetectLocaleNoneSet(t *testing.T) {
	LoadBundle()
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "")
	t.Setenv("LANG", "")
	DetectLocale()
	if Lang() != "en" {
		t.Errorf("expected lang 'en', got %q", Lang())
	}
}

func TestDetectLocaleC(t *testing.T) {
	LoadBundle()
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "")
	t.Setenv("LANG", "C")
	DetectLocale()
	if Lang() != "en" {
		t.Errorf("expected lang 'en', got %q", Lang())
	}
}

func TestDetectLocalePOSIX(t *testing.T) {
	LoadBundle()
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "")
	t.Setenv("LANG", "POSIX")
	DetectLocale()
	if Lang() != "en" {
		t.Errorf("expected lang 'en', got %q", Lang())
	}
}

func TestDetectLocaleBareCode(t *testing.T) {
	LoadBundle()
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "")
	t.Setenv("LANG", "uk")
	DetectLocale()
	if Lang() != "uk" {
		t.Errorf("expected lang 'uk', got %q", Lang())
	}
}

func TestDetectLocaleMalformedFallsBack(t *testing.T) {
	LoadBundle()
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "")
	t.Setenv("LANG", "zz_ZZ.UTF-8")
	DetectLocale()
	if Lang() != "en" {
		t.Errorf("expected lang 'en' for unsupported locale, got %q", Lang())
	}
}

// T024: Verify unsupported locale de_DE.UTF-8 resolves to en.
func TestDetectLocaleGermanFallsBackToEn(t *testing.T) {
	LoadBundle()
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "")
	t.Setenv("LANG", "de_DE.UTF-8")
	DetectLocale()
	if Lang() != "en" {
		t.Errorf("expected lang 'en' for German (unsupported), got %q", Lang())
	}
}

// T030: Verify LoadBundle auto-discovers all .toml files.
func TestLoadBundleAutoDiscovery(t *testing.T) {
	LoadBundle()
	// Count should equal the number of .toml files in locales/
	// Currently: en.toml, uk.toml = 2
	if BundleLen() != 2 {
		t.Errorf("expected 2 bundles, got %d", BundleLen())
	}
}

func TestTranslateEnglish(t *testing.T) {
	LoadBundle()
	SetLang("en")
	got := Translate("root.short")
	want := "mytets — a lightweight CLI tool"
	if got != want {
		t.Errorf("Translate(root.short) = %q, want %q", got, want)
	}
}

func TestTranslateUkrainian(t *testing.T) {
	LoadBundle()
	SetLang("uk")
	got := Translate("root.short")
	want := "mytets — легкий інструмент командного рядка"
	if got != want {
		t.Errorf("Translate(root.short) = %q, want %q", got, want)
	}
}

func TestTranslateFallbackToEnglish(t *testing.T) {
	LoadBundle()
	SetLang("uk")
	// If uk.toml is missing a key but en.toml has it, en value is returned.
	// Since both files have the same keys, test with a synthetic scenario
	// by removing a key from uk bundle.
	delete(bundle["uk"], "root.short")
	defer func() {
		// Reload to restore state.
		LoadBundle()
	}()
	got := Translate("root.short")
	want := "mytets — a lightweight CLI tool"
	if got != want {
		t.Errorf("fallback: Translate(root.short) = %q, want %q", got, want)
	}
}

func TestTranslateFallbackToKeyName(t *testing.T) {
	LoadBundle()
	SetLang("en")
	got := Translate("nonexistent.key")
	if got != "nonexistent.key" {
		t.Errorf("expected key name fallback, got %q", got)
	}
}

func TestSetLang(t *testing.T) {
	LoadBundle()
	SetLang("uk")
	if Lang() != "uk" {
		t.Errorf("SetLang(uk): Lang() = %q", Lang())
	}
	SetLang("en")
	if Lang() != "en" {
		t.Errorf("SetLang(en): Lang() = %q", Lang())
	}
}

func TestExtractLangCode(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"uk_UA.UTF-8", "uk"},
		{"en_US.UTF-8", "en"},
		{"de_DE", "de"},
		{"fr.UTF-8", "fr"},
		{"uk", "uk"},
		{"", "en"},
		{"C", "en"},
		{"POSIX", "en"},
		{"  ", "en"},
	}
	for _, tc := range tests {
		got := extractLangCode(tc.input)
		if got != tc.want {
			t.Errorf("extractLangCode(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
