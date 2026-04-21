package one

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/igorzel/mytets/internal/flags"
	"github.com/igorzel/mytets/internal/phrases"
	"github.com/spf13/cobra"
)

// T005: TestNew - Verify command is properly constructed
func TestNew(t *testing.T) {
	cfg := flags.ParserConfig{}
	cmd := New(cfg)
	if cmd == nil {
		t.Fatal("New() returned nil command")
	}
	if cmd.Use != "one" {
		t.Errorf("Command Use = %q, want %q", cmd.Use, "one")
	}
	if cmd.Short == "" {
		t.Error("Command Short description is empty")
	}
	if cmd.Short != "Display one random phrase" {
		t.Errorf("Command Short = %q, want %q", cmd.Short, "Display one random phrase")
	}
	if cmd.RunE == nil {
		t.Error("Command RunE is nil")
	}
}

// T004: TestOutputPlain - Verify plain text output
func TestOutputPlain(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
	}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := outputPlain(cmd)
	if err != nil {
		t.Fatalf("outputPlain() error = %v", err)
	}

	output := strings.TrimSpace(buf.String())
	if output == "" {
		t.Fatal("outputPlain() produced empty output")
	}
	if !contains(phrases.Messages(), output) {
		t.Errorf("outputPlain() output = %q, not found in embedded phrases", output)
	}
}

// T013: TestOutputJSON - Verify JSON output is valid and compact
func TestOutputJSON(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
	}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := outputJSON(cmd)
	if err != nil {
		t.Fatalf("outputJSON() error = %v", err)
	}

	output := strings.TrimSpace(buf.String())

	// Verify output can be parsed as JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("outputJSON() output is not valid JSON: %v", err)
	}

	// Verify message field exists and has correct value
	msgVal, ok := parsed["message"]
	if !ok {
		t.Error("JSON output missing 'message' field")
	}
	msg, ok := msgVal.(string)
	if !ok {
		t.Fatalf("message field type = %T, want string", msgVal)
	}
	if !contains(phrases.Messages(), msg) {
		t.Errorf("message field = %q, not found in embedded phrases", msg)
	}
}

// T014: TestResponseJSONFormat - Verify exact JSON format with no pretty-printing
func TestResponseJSONFormat(t *testing.T) {
	resp := Response{Message: "example"}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	jsonStr := string(data)
	expected := `{"message":"example"}`

	if jsonStr != expected {
		t.Errorf("JSON format = %q, want %q", jsonStr, expected)
	}

	// Verify compact format (no newlines or excessive indentation)
	if bytes.Contains(data, []byte("\n")) || bytes.Contains(data, []byte("  ")) {
		t.Error("JSON output is not compact (contains newlines or indentation)")
	}
}

// Helper: TestResponseStruct - Verify Response struct has correct field
func TestResponseStruct(t *testing.T) {
	resp := Response{Message: "test"}
	if resp.Message != "test" {
		t.Errorf("Response.Message = %q, want %q", resp.Message, "test")
	}
}

func TestOutputPlainPropagatesPhraseError(t *testing.T) {
	original := randomMessage
	t.Cleanup(func() { randomMessage = original })
	randomMessage = func() (string, error) {
		return "", errors.New("boom")
	}

	cmd := &cobra.Command{Use: "test"}
	err := outputPlain(cmd)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to select phrase") {
		t.Fatalf("unexpected error: %v", err)
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
