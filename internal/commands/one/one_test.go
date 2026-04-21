package one

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/igorzel/mytets/internal/flags"
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

	output := buf.String()
	expected := message + "\n"
	if output != expected {
		t.Errorf("outputPlain() output = %q, want %q", output, expected)
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

	output := buf.String()

	// Verify output can be parsed as JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(output[:len(output)-1]), &parsed); err != nil {
		t.Fatalf("outputJSON() output is not valid JSON: %v", err)
	}

	// Verify message field exists and has correct value
	msgVal, ok := parsed["message"]
	if !ok {
		t.Error("JSON output missing 'message' field")
	}
	if msgVal != message {
		t.Errorf("message field = %q, want %q", msgVal, message)
	}
}

// T014: TestResponseJSONFormat - Verify exact JSON format with no pretty-printing
func TestResponseJSONFormat(t *testing.T) {
	resp := Response{Message: message}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	jsonStr := string(data)
	expected := `{"message":"Fake message, tbd"}`

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
