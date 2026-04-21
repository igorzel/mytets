package cli_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/igorzel/mytets/internal/cli"
	"github.com/igorzel/mytets/internal/version"
)

func TestExecuteArgsVersionPlain(t *testing.T) {
	original := version.Version
	t.Cleanup(func() { version.Version = original })
	version.Version = "1.0.0"

	stdout, stderr, code := cli.ExecuteArgs([]string{"version"})

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
	if stderr != "" {
		t.Errorf("expected empty stderr, got %q", stderr)
	}
	if strings.TrimRight(stdout, "\n") != "1.0.0" {
		t.Errorf("expected stdout '1.0.0', got %q", stdout)
	}
}

func TestExecuteArgsVersionJSON(t *testing.T) {
	original := version.Version
	t.Cleanup(func() { version.Version = original })
	version.Version = "2.3.4"

	stdout, stderr, code := cli.ExecuteArgs([]string{"version", "--output", "json"})

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
	if stderr != "" {
		t.Errorf("expected empty stderr, got %q", stderr)
	}
	var payload map[string]string
	if err := json.Unmarshal([]byte(strings.TrimRight(stdout, "\n")), &payload); err != nil {
		t.Fatalf("stdout not valid JSON: %v; stdout=%q", err, stdout)
	}
	if payload["version"] != "2.3.4" {
		t.Errorf("expected version '2.3.4', got %q", payload["version"])
	}
}

func TestExecuteArgsVersionShortFlag(t *testing.T) {
	original := version.Version
	t.Cleanup(func() { version.Version = original })
	version.Version = "5.6.7"

	stdout, stderr, code := cli.ExecuteArgs([]string{"version", "-o", "json"})

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
	if stderr != "" {
		t.Errorf("expected empty stderr, got %q", stderr)
	}
	var payload map[string]string
	if err := json.Unmarshal([]byte(strings.TrimRight(stdout, "\n")), &payload); err != nil {
		t.Fatalf("stdout not valid JSON: %v; stdout=%q", err, stdout)
	}
	if payload["version"] != "5.6.7" {
		t.Errorf("expected version '5.6.7', got %q", payload["version"])
	}
}

func TestExecuteArgsVersionUnsupportedFormat(t *testing.T) {
	_, stderr, code := cli.ExecuteArgs([]string{"version", "--output", "yaml"})

	if code == 0 {
		t.Error("expected non-zero exit code for unsupported format, got 0")
	}
	if stderr == "" {
		t.Error("expected non-empty stderr for unsupported format")
	}
}

func TestExecuteArgsVersionDevFallback(t *testing.T) {
	original := version.Version
	t.Cleanup(func() { version.Version = original })
	version.Version = "dev"

	stdout, stderr, code := cli.ExecuteArgs([]string{"version"})

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
	if stderr != "" {
		t.Errorf("expected empty stderr, got %q", stderr)
	}
	if strings.TrimRight(stdout, "\n") != "dev" {
		t.Errorf("expected stdout 'dev', got %q", stdout)
	}
}

func TestExecuteArgsUnknownSubcommand(t *testing.T) {
	_, _, code := cli.ExecuteArgs([]string{"nonexistent"})
	if code == 0 {
		t.Error("expected non-zero exit code for unknown subcommand, got 0")
	}
}

func TestExecuteArgsVersionUnknownFlag(t *testing.T) {
	_, stderr, code := cli.ExecuteArgs([]string{"version", "--foo"})
	if code == 0 {
		t.Error("expected non-zero exit code for unknown flag, got 0")
	}
	if stderr == "" {
		t.Error("expected non-empty stderr for unknown flag")
	}
}

func TestExecuteArgsVersionExtraArgs(t *testing.T) {
	_, stderr, code := cli.ExecuteArgs([]string{"version", "unexpected"})
	if code == 0 {
		t.Error("expected non-zero exit for extra positional arg, got 0")
	}
	if stderr == "" {
		t.Error("expected non-empty stderr for extra positional arg")
	}
}
