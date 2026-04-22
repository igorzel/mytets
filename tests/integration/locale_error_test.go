package integration

import (
	"strings"
	"testing"
)

func TestLocaleErrorInvalidFormatUkrainian(t *testing.T) {
	bin := buildBinary(t, "")

	_, stderr, code := runBinaryWithEnv(t, bin, []string{"LANG=uk_UA.UTF-8"}, "--output", "xml", "one")

	if code == 0 {
		t.Error("expected non-zero exit code")
	}
	want := `непідтримуваний формат виводу: "xml"`
	if !strings.Contains(stderr, want) {
		t.Errorf("stderr missing %q\nGot: %q", want, stderr)
	}
}

func TestLocaleErrorUnknownCommandUkrainian(t *testing.T) {
	bin := buildBinary(t, "")

	_, stderr, code := runBinaryWithEnv(t, bin, []string{"LANG=uk_UA.UTF-8"}, "foo")

	if code == 0 {
		t.Error("expected non-zero exit code")
	}
	want := `невідома команда "foo" для "mytets"`
	if !strings.Contains(stderr, want) {
		t.Errorf("stderr missing %q\nGot: %q", want, stderr)
	}
}

func TestLocaleErrorUnknownFlagUkrainian(t *testing.T) {
	bin := buildBinary(t, "")

	_, stderr, code := runBinaryWithEnv(t, bin, []string{"LANG=uk_UA.UTF-8"}, "version", "--foo")

	if code == 0 {
		t.Error("expected non-zero exit code")
	}
	want := "невідомий прапор: --foo"
	if !strings.Contains(stderr, want) {
		t.Errorf("stderr missing %q\nGot: %q", want, stderr)
	}
}
