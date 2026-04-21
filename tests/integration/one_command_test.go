package integration

import (
"bytes"
"encoding/json"
"testing"

"github.com/igorzel/mytets/internal/cli"
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

expected := "Fake message, tbd\n"
if stdout != expected {
t.Errorf("stdout = %q, want %q", stdout, expected)
}
}

// T016: TestOneCommandJSON - Verify JSON invocation via CLI
func TestOneCommandJSON(t *testing.T) {
stdout, stderr, exitCode := cli.ExecuteArgs([]string{"one", "--output", "json"})

if exitCode != 0 {
t.Errorf("Exit code = %d, want 0", exitCode)
}

if stderr != "" {
t.Errorf("stderr should be empty, got %q", stderr)
}

// Verify output is valid JSON
output := bytes.TrimSpace([]byte(stdout))
var parsed map[string]interface{}
if err := json.Unmarshal(output, &parsed); err != nil {
t.Fatalf("stdout is not valid JSON: %v", err)
}

// Verify message field
msgVal, ok := parsed["message"]
if !ok {
t.Error("JSON output missing 'message' field")
}

if msgVal != "Fake message, tbd" {
t.Errorf("message field = %q, want %q", msgVal, "Fake message, tbd")
}
}

// T017: TestOneCommandJSONExitCode - Verify exit code in JSON mode
func TestOneCommandJSONExitCode(t *testing.T) {
_, _, exitCode := cli.ExecuteArgs([]string{"one", "-o", "json"})
if exitCode != 0 {
t.Errorf("Exit code in JSON mode = %d, want 0", exitCode)
}
}

// Test: TestOneCommandPlainAlternativeFlag - Verify plain text with explicit text flag
func TestOneCommandPlainAlternativeFlag(t *testing.T) {
stdout, stderr, exitCode := cli.ExecuteArgs([]string{"one", "--output", "text"})

if exitCode != 0 {
t.Errorf("Exit code = %d, want 0", exitCode)
}

if stderr != "" {
t.Errorf("stderr should be empty, got %q", stderr)
}

expected := "Fake message, tbd\n"
if stdout != expected {
t.Errorf("stdout = %q, want %q", stdout, expected)
}
}

// Test: TestOneCommandInvalidFormat - Verify error handling for invalid format
func TestOneCommandInvalidFormat(t *testing.T) {
_, stderr, exitCode := cli.ExecuteArgs([]string{"one", "--output", "invalid"})

if exitCode == 0 {
t.Error("Exit code should be non-zero for invalid format")
}

if stderr == "" {
t.Error("stderr should contain error message")
}
}
