// Package integration contains end-to-end tests that build the real binary and
// invoke it as a child process, verifying stdout, stderr and exit code.
package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"
)

// buildBinary compiles the binary into a temp directory and returns the path.
// An optional ldflags string can be supplied to inject version metadata.
func buildBinary(t *testing.T, ldflags string) string {
	t.Helper()

	dir := t.TempDir()
	name := "mytets"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	binPath := filepath.Join(dir, name)

	args := []string{"build", "-o", binPath}
	if ldflags != "" {
		args = append(args, "-ldflags", ldflags)
	}
	// Resolve module root relative to this test file.
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	moduleRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")

	args = append(args, "./cmd/mytets")
	cmd := exec.Command("go", args...)
	cmd.Dir = moduleRoot

	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, errBuf.String())
	}
	return binPath
}

// runBinary invokes the binary with args and returns stdout, stderr and the
// exit code.
func runBinary(t *testing.T, binPath string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command(binPath, args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("unexpected error running binary: %v", err)
		}
	}
	return outBuf.String(), errBuf.String(), exitCode
}

// ── US1: Developer Queries Application Version ──────────────────────────────

// T012 — basic invocation: exit code 0, single-line stdout, empty stderr.
func TestVersionCommandBasic(t *testing.T) {
	bin := buildBinary(t, "")

	stdout, stderr, code := runBinary(t, bin, "version")

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
	if stderr != "" {
		t.Errorf("expected empty stderr, got %q", stderr)
	}
	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	if len(lines) != 1 {
		t.Errorf("expected exactly 1 output line, got %d: %q", len(lines), stdout)
	}
}

// T012 — fallback to "dev" when no ldflags are supplied.
func TestVersionCommandDevFallback(t *testing.T) {
	bin := buildBinary(t, "")

	stdout, stderr, code := runBinary(t, bin, "version")

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
	if stderr != "" {
		t.Errorf("expected empty stderr, got %q", stderr)
	}
	got := strings.TrimRight(stdout, "\n")
	if got != "dev" {
		t.Errorf("expected fallback 'dev', got %q", got)
	}
}

// T038 — ldflags-injected version is printed exactly.
func TestVersionCommandLdflagsInjected(t *testing.T) {
	const injected = "1.2.3"
	ldflags := fmt.Sprintf("-X github.com/igorzel/mytets/internal/version.Version=%s", injected)
	bin := buildBinary(t, ldflags)

	stdout, stderr, code := runBinary(t, bin, "version")

	if code != 0 {
		t.Errorf("expected exit code 0, got %d; stderr: %q", code, stderr)
	}
	if stderr != "" {
		t.Errorf("expected empty stderr, got %q", stderr)
	}
	got := strings.TrimRight(stdout, "\n")
	if got != injected {
		t.Errorf("expected %q, got %q", injected, got)
	}
}

// T038 — second ldflags variant to match spec scenario 2.
func TestVersionCommandLdflagsVariant(t *testing.T) {
	const injected = "0.9.0"
	ldflags := fmt.Sprintf("-X github.com/igorzel/mytets/internal/version.Version=%s", injected)
	bin := buildBinary(t, ldflags)

	stdout, _, code := runBinary(t, bin, "version")

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
	got := strings.TrimRight(stdout, "\n")
	if got != injected {
		t.Errorf("expected %q, got %q", injected, got)
	}
}

// ── US2: Scripted / CI Version Check ────────────────────────────────────────

// T017 — regex-compatible output and empty stderr (plain text, injected semver).
func TestVersionCommandRegexCompatible(t *testing.T) {
	const injected = "2.0.0"
	ldflags := fmt.Sprintf("-X github.com/igorzel/mytets/internal/version.Version=%s", injected)
	bin := buildBinary(t, ldflags)

	stdout, stderr, code := runBinary(t, bin, "version")

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
	if stderr != "" {
		t.Errorf("expected empty stderr, got %q", stderr)
	}
	got := strings.TrimRight(stdout, "\n")
	semverRe := regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
	if !semverRe.MatchString(got) {
		t.Errorf("output %q does not match semver pattern", got)
	}
}

// T018 — dev fallback is output when no version is injected.
func TestVersionCommandDevFallbackScripting(t *testing.T) {
	bin := buildBinary(t, "")

	stdout, _, code := runBinary(t, bin, "version")

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
	got := strings.TrimRight(stdout, "\n")
	if got != "dev" {
		t.Errorf("expected fallback 'dev', got %q", got)
	}
}

// T019 — unknown flag yields non-zero exit code and non-empty stderr.
func TestVersionCommandUnknownFlagNonZeroExit(t *testing.T) {
	bin := buildBinary(t, "")

	_, stderr, code := runBinary(t, bin, "version", "--foo")

	if code == 0 {
		t.Error("expected non-zero exit code for unknown flag, got 0")
	}
	if stderr == "" {
		t.Error("expected non-empty stderr for unknown flag")
	}
}

// T028 — plain text output (no flag) versus JSON output (--output json).
func TestVersionCommandPlainOutput(t *testing.T) {
	const injected = "3.1.4"
	ldflags := fmt.Sprintf("-X github.com/igorzel/mytets/internal/version.Version=%s", injected)
	bin := buildBinary(t, ldflags)

	stdout, stderr, code := runBinary(t, bin, "version")

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
	if stderr != "" {
		t.Errorf("expected empty stderr, got %q", stderr)
	}
	got := strings.TrimRight(stdout, "\n")
	if got != injected {
		t.Errorf("plain output: expected %q, got %q", injected, got)
	}
}

func TestVersionCommandJSONOutput(t *testing.T) {
	const injected = "3.1.4"
	ldflags := fmt.Sprintf("-X github.com/igorzel/mytets/internal/version.Version=%s", injected)
	bin := buildBinary(t, ldflags)

	stdout, stderr, code := runBinary(t, bin, "version", "--output", "json")

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
	if stderr != "" {
		t.Errorf("expected empty stderr, got %q", stderr)
	}

	var payload map[string]string
	if err := json.Unmarshal([]byte(strings.TrimRight(stdout, "\n")), &payload); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\nstdout: %q", err, stdout)
	}
	if payload["version"] != injected {
		t.Errorf("JSON version: expected %q, got %q", injected, payload["version"])
	}
}

// T039 — unsupported output format returns clear stderr and non-zero exit.
func TestVersionCommandUnsupportedOutputFormat(t *testing.T) {
	bin := buildBinary(t, "")

	_, stderr, code := runBinary(t, bin, "version", "--output", "yaml")

	if code == 0 {
		t.Error("expected non-zero exit code for unsupported output format, got 0")
	}
	if stderr == "" {
		t.Error("expected non-empty stderr for unsupported output format")
	}
}

// T037 — --help contains expected command and flag descriptions.
func TestVersionCommandHelpText(t *testing.T) {
	bin := buildBinary(t, "")

	stdout, _, code := runBinary(t, bin, "version", "--help")

	if code != 0 {
		t.Errorf("expected exit code 0 for --help, got %d", code)
	}
	if !strings.Contains(stdout, "version") {
		t.Errorf("help text missing 'version': %q", stdout)
	}
	if !strings.Contains(stdout, "--output") {
		t.Errorf("help text missing '--output' flag: %q", stdout)
	}
}

// T029 — performance: version command completes in under 100 ms.
func TestVersionCommandPerformance(t *testing.T) {
	bin := buildBinary(t, "")

	start := time.Now()
	_, _, code := runBinary(t, bin, "version")
	elapsed := time.Since(start)

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
	const threshold = 100 * time.Millisecond
	if elapsed > threshold {
		t.Errorf("version command took %v, want < %v", elapsed, threshold)
	}
}

// T028 – short flag alias -o works for JSON output.
func TestVersionCommandShortOutputFlag(t *testing.T) {
	const injected = "5.0.0"
	ldflags := fmt.Sprintf("-X github.com/igorzel/mytets/internal/version.Version=%s", injected)
	bin := buildBinary(t, ldflags)

	stdout, stderr, code := runBinary(t, bin, "version", "-o", "json")

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
	if stderr != "" {
		t.Errorf("expected empty stderr, got %q", stderr)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimRight(stdout, "\n")), &payload); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\nstdout: %q", err, stdout)
	}
	if v, ok := payload["version"].(string); !ok || v != injected {
		t.Errorf("JSON version: expected %q, got %v", injected, payload["version"])
	}
}

// Ensure the test file can be compiled even without test functions above.
var _ = os.DevNull
