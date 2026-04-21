package integration

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/igorzel/mytets/internal/cli"
	"github.com/igorzel/mytets/internal/phrases"
)

// T007: TestOneCommandPlain - Verify plain text invocation via CLI
func TestOneCommandPlain(t *testing.T) {
	stdout, stderr, exitCode := cli.ExecuteArgs([]string{"one"})

	if exitCode != 0 {
		t.Errorf("Exit code = %d, want 0", exitCode)
	}

	if stderr != "" {
		t.Errorf("stderr should be empty, got %q", stderr)
	}

	plain := strings.TrimSpace(stdout)
	if plain == "" {
		t.Fatal("stdout should not be empty")
	}
	if !contains(phrases.Messages(), plain) {
		t.Errorf("stdout phrase %q is not in embedded phrase set", plain)
	}
}

// T016: TestOneCommandJSON - Verify JSON invocation via CLI
func TestOneCommandJSON(t *testing.T) {
	stdout, stderr, exitCode := cli.ExecuteArgs([]string{"--output", "json", "one"})

	if exitCode != 0 {
		t.Errorf("Exit code = %d, want 0", exitCode)
	}

	if stderr != "" {
		t.Errorf("stderr should be empty, got %q", stderr)
	}

	// Verify output is valid JSON
	output := strings.TrimSpace(stdout)
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("stdout is not valid JSON: %v", err)
	}

	// Verify message field
	msgVal, ok := parsed["message"]
	if !ok {
		t.Error("JSON output missing 'message' field")
	}

	msg, ok := msgVal.(string)
	if !ok {
		t.Fatalf("message field type = %T, want string", msgVal)
	}
	if !contains(phrases.Messages(), msg) {
		t.Errorf("message field %q is not in embedded phrase set", msg)
	}
}

// T017: TestOneCommandJSONExitCode - Verify exit code in JSON mode
func TestOneCommandJSONExitCode(t *testing.T) {
	_, _, exitCode := cli.ExecuteArgs([]string{"--output", "json", "one"})
	if exitCode != 0 {
		t.Errorf("Exit code in JSON mode = %d, want 0", exitCode)
	}
}

// Test: TestOneCommandPlainAlternativeFlag - Verify plain text with explicit text flag
func TestOneCommandPlainAlternativeFlag(t *testing.T) {
	stdout, stderr, exitCode := cli.ExecuteArgs([]string{"--output", "text", "one"})

	if exitCode != 0 {
		t.Errorf("Exit code = %d, want 0", exitCode)
	}

	if stderr != "" {
		t.Errorf("stderr should be empty, got %q", stderr)
	}

	plain := strings.TrimSpace(stdout)
	if plain == "" {
		t.Fatal("stdout should not be empty")
	}
	if !contains(phrases.Messages(), plain) {
		t.Errorf("stdout phrase %q is not in embedded phrase set", plain)
	}
}

// Test: TestOneCommandInvalidFormat - Verify error handling for invalid format
func TestOneCommandInvalidFormat(t *testing.T) {
	_, stderr, exitCode := cli.ExecuteArgs([]string{"--output", "invalid", "one"})

	if exitCode == 0 {
		t.Error("Exit code should be non-zero for invalid format")
	}

	if stderr == "" {
		t.Error("stderr should contain error message")
	}
}

func TestOneCommandRandomnessAcrossRuns(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 40; i++ {
		stdout, stderr, exitCode := cli.ExecuteArgs([]string{"one"})
		if exitCode != 0 {
			t.Fatalf("run %d: exit code = %d, want 0 (stderr=%q)", i, exitCode, stderr)
		}
		seen[strings.TrimSpace(stdout)] = true
	}
	if len(seen) < 2 {
		t.Fatalf("expected at least 2 different phrases across runs, got %d", len(seen))
	}
}

func contains(values []string, target string) bool {
	for _, v := range values {
		if v == target {
			return true
		}
	}
	return false
}
