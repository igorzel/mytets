package list

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/igorzel/mytets/internal/flags"
	"github.com/spf13/cobra"
)

func TestNew(t *testing.T) {
	cfg := flags.ParserConfig{}
	cmd := New(cfg)
	if cmd == nil {
		t.Fatal("New() returned nil command")
	}
	if cmd.Use != "list" {
		t.Errorf("Command Use = %q, want %q", cmd.Use, "list")
	}
	if cmd.Short == "" {
		t.Error("Command Short description is empty")
	}
	if cmd.RunE == nil {
		t.Error("Command RunE is nil")
	}
}

func TestCountFlagDefault(t *testing.T) {
	cfg := flags.ParserConfig{}
	cmd := New(cfg)

	flag := cmd.Flags().Lookup("count")
	if flag == nil {
		t.Fatal("--count flag not found")
	}
	if flag.DefValue != "5" {
		t.Errorf("--count default = %q, want %q", flag.DefValue, "5")
	}
}

func TestCountFlagIsCommandSpecific(t *testing.T) {
	cfg := flags.ParserConfig{}
	cmd := New(cfg)

	// The flag should be on the command, not inherited
	if cmd.Flags().Lookup("count") == nil {
		t.Error("--count should be a local flag on the list command")
	}
	if cmd.InheritedFlags().Lookup("count") != nil {
		t.Error("--count should not be an inherited flag")
	}
}

func TestOutputPlainFormat(t *testing.T) {
	original := messageSource
	defer func() { messageSource = original }()
	messageSource = func() []string {
		return []string{"alpha", "bravo", "charlie", "delta", "echo"}
	}

	cfg := flags.ParserConfig{Output: flags.OutputFormatText}
	cmd := New(cfg)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--count", "3"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	lines := nonEmptyLines(buf.String())
	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3", len(lines))
	}
	assertUniqueLines(t, lines)
}

func TestOutputPlainDefaultCount(t *testing.T) {
	original := messageSource
	defer func() { messageSource = original }()
	messageSource = func() []string {
		return []string{"a", "b", "c", "d", "e", "f", "g"}
	}

	cfg := flags.ParserConfig{Output: flags.OutputFormatText}
	cmd := New(cfg)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	lines := nonEmptyLines(buf.String())
	if len(lines) != 5 {
		t.Fatalf("got %d lines, want 5 (default)", len(lines))
	}
}

func TestErrorWhenNoPhrasesAvailable(t *testing.T) {
	original := messageSource
	defer func() { messageSource = original }()
	messageSource = func() []string { return nil }

	cfg := flags.ParserConfig{Output: flags.OutputFormatText}
	cmd := New(cfg)

	var outBuf, errBuf bytes.Buffer
	cmd.SetOut(&outBuf)
	cmd.SetErr(&errBuf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() should return an error when no phrases available")
	}
	if !strings.Contains(err.Error(), "no phrases available") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "no phrases available")
	}
}

func TestErrorWhenEmptySliceReturned(t *testing.T) {
	original := messageSource
	defer func() { messageSource = original }()
	messageSource = func() []string { return []string{} }

	cfg := flags.ParserConfig{Output: flags.OutputFormatText}
	cmd := New(cfg)

	var outBuf, errBuf bytes.Buffer
	cmd.SetOut(&outBuf)
	cmd.SetErr(&errBuf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() should return an error when phrase source returns empty slice")
	}
}

// --- JSON output unit tests (T010) ---

func TestOutputJSONFormat(t *testing.T) {
	original := messageSource
	defer func() { messageSource = original }()
	messageSource = func() []string {
		return []string{"alpha", "bravo", "charlie"}
	}

	cfg := flags.ParserConfig{Output: flags.OutputFormatJSON}
	cmd := New(cfg)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--count", "2"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	trimmed := strings.TrimSpace(buf.String())

	var items []Response
	if err := json.Unmarshal([]byte(trimmed), &items); err != nil {
		t.Fatalf("output is not valid JSON array: %v\nraw: %q", err, trimmed)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}
	for i, item := range items {
		if item.Message == "" {
			t.Errorf("item[%d].Message is empty", i)
		}
	}
}

func TestOutputJSONCompactness(t *testing.T) {
	original := messageSource
	defer func() { messageSource = original }()
	messageSource = func() []string { return []string{"x"} }

	cfg := flags.ParserConfig{Output: flags.OutputFormatJSON}
	cmd := New(cfg)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--count", "1"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	trimmed := strings.TrimSpace(buf.String())
	expected := `[{"message":"x"}]`
	if trimmed != expected {
		t.Errorf("JSON output = %q, want %q", trimmed, expected)
	}
}

func TestOutputJSONArrayStructure(t *testing.T) {
	original := messageSource
	defer func() { messageSource = original }()
	messageSource = func() []string {
		return []string{"a", "b", "c", "d", "e"}
	}

	cfg := flags.ParserConfig{Output: flags.OutputFormatJSON}
	cmd := New(cfg)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--count", "3"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	trimmed := strings.TrimSpace(buf.String())

	// Must start with [ and end with ]
	if !strings.HasPrefix(trimmed, "[") || !strings.HasSuffix(trimmed, "]") {
		t.Errorf("JSON output does not look like an array: %q", trimmed)
	}

	// Parse and verify field names
	var raw []map[string]interface{}
	if err := json.Unmarshal([]byte(trimmed), &raw); err != nil {
		t.Fatalf("failed to parse as array of objects: %v", err)
	}
	for i, obj := range raw {
		if _, ok := obj["message"]; !ok {
			t.Errorf("item[%d] missing 'message' field", i)
		}
		if len(obj) != 1 {
			t.Errorf("item[%d] has %d fields, want 1", i, len(obj))
		}
	}
}

// --- helpers ---

func nonEmptyLines(s string) []string {
	var lines []string
	for _, line := range strings.Split(s, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	return lines
}

func assertUniqueLines(t *testing.T, lines []string) {
	t.Helper()
	seen := make(map[string]struct{}, len(lines))
	for _, line := range lines {
		if _, exists := seen[line]; exists {
			t.Errorf("duplicate line found: %q", line)
		}
		seen[line] = struct{}{}
	}
}

// suppress unused lint for cobra import in tests
var _ = (*cobra.Command)(nil)
