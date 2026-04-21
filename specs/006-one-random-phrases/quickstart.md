# Quick Start: One Command - Random Phrase Feature

**Feature**: Implement `mytets one` command  
**Audience**: Developers implementing and testing this feature  
**Duration**: ~2-3 hours for initial implementation + testing  

---

## 1. Prerequisites

### System Requirements
- Go 1.25+ (project uses 1.26.2)
- POSIX shell (for build scripts)
- `go build` and `go test` commands available

### Project Setup
```bash
cd /home/igor/dev/workspace/mytets
go mod download  # Download dependencies (github.com/spf13/cobra already present)
```

---

## 2. Create Embedded Phrases File

Create the embedded data file that will be compiled into the binary:

**File**: `internal/phrases/phrases.json`

```bash
mkdir -p internal/phrases
```

**Content** (example):
```json
{
    "messages": [
        {
            "text": "Fake message, tbd"
        },
        {
            "text": "Example message 1"
        },
        {
            "text": "Example message 2"
        },
        {
            "text": "Example message 3"
        }
    ]
}
```

**Note**: This file is embedded at compile-time via `//go:embed` directive. Must be valid JSON with non-empty `messages` array.

---

## 3. Create Phrase Loading Package

Create a new internal package for phrase loading and random selection:

**File**: `internal/phrases/phrases.go`

```go
package phrases

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

//go:embed phrases.json
var phrasesJSON string

type data struct {
	Messages []struct {
		Text string `json:"text"`
	} `json:"messages"`
}

var phrases data
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func init() {
	if err := json.Unmarshal([]byte(phrasesJSON), &phrases); err != nil {
		panic(fmt.Sprintf("failed to parse embedded phrases.json: %v", err))
	}
	if len(phrases.Messages) == 0 {
		panic("phrases.json contains no messages")
	}
}

// GetMessages returns all available phrase messages
func GetMessages() []string {
	result := make([]string, len(phrases.Messages))
	for i, m := range phrases.Messages {
		result[i] = m.Text
	}
	return result
}

// GetRandomPhrase returns a single random phrase
func GetRandomPhrase() string {
	messages := GetMessages()
	if len(messages) == 0 {
		return "" // Caught by init() validation
	}
	return messages[rng.Intn(len(messages))]
}
```

**File**: `internal/phrases/phrases_test.go`

```go
package phrases

import (
	"testing"
)

func TestGetMessages(t *testing.T) {
	messages := GetMessages()
	if len(messages) == 0 {
		t.Fatal("expected non-empty messages")
	}
	for _, msg := range messages {
		if msg == "" {
			t.Error("expected non-empty phrase text")
		}
	}
}

func TestGetRandomPhrase(t *testing.T) {
	phrase := GetRandomPhrase()
	if phrase == "" {
		t.Error("expected non-empty random phrase")
	}

	// Verify it's in the list of valid phrases
	messages := GetMessages()
	found := false
	for _, m := range messages {
		if m == phrase {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("GetRandomPhrase returned invalid phrase: %q", phrase)
	}
}

func TestRandomness(t *testing.T) {
	results := make(map[string]int)
	for i := 0; i < 100; i++ {
		phrase := GetRandomPhrase()
		results[phrase]++
	}

	if len(results) == 1 {
		t.Error("expected multiple different phrases in 100 runs")
	}
}
```

---

## 4. Create One Command

Create the command implementation:

**File**: `internal/commands/one/one.go`

```go
package one

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/igorzel/mytets/internal/flags"
	"github.com/igorzel/mytets/internal/phrases"
)

func New(cfg flags.ParserConfig) *cobra.Command {
	var outputRaw string

	cmd := &cobra.Command{
		Use:   "one",
		Short: "Display a random phrase",
		Long:  "The one command outputs a random phrase in plain text or JSON format.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			format, err := flags.ParseOutputFormat(outputRaw)
			if err != nil {
				return fmt.Errorf("invalid output format: %w", err)
			}

			switch format {
			case flags.OutputFormatJSON:
				return outputJSON(cmd)
			default:
				return outputPlain(cmd)
			}
		},
	}

	cmd.Flags().StringVarP(
		&outputRaw,
		"output", "o",
		string(cfg.Output),
		`Output format: "text" (default) or "json"`,
	)

	return cmd
}

func outputPlain(cmd *cobra.Command) error {
	phrase := phrases.GetRandomPhrase()
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), phrase)
	return nil
}

func outputJSON(cmd *cobra.Command) error {
	phrase := phrases.GetRandomPhrase()
	output, err := json.Marshal(map[string]string{"message": phrase})
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	_, _ = fmt.Fprint(cmd.OutOrStdout(), string(output))
	fmt.Fprintln(cmd.OutOrStdout()) // Newline after JSON
	return nil
}
```

**File**: `internal/commands/one/one_test.go`

```go
package one

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/spf13/cobra"
	"github.com/igorzel/mytets/internal/flags"
)

func TestOneCommand(t *testing.T) {
	cmd := New(flags.ParserConfig{Output: "text"})

	// Test plain text output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("RunE failed: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("expected non-empty output")
	}
}

func TestOneCommandJSON(t *testing.T) {
	cmd := New(flags.ParserConfig{Output: "json"})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	
	// Simulate --output json flag
	cmd.Flags().Set("output", "json")
	
	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("RunE failed: %v", err)
	}

	output := buf.String()
	
	// Verify JSON is valid
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if _, ok := parsed["message"]; !ok {
		t.Error("JSON output missing 'message' field")
	}
}
```

---

## 5. Register Command in Root CLI

Update the root CLI to register the new command:

**File**: `internal/cli/root.go` (modify existing file)

Add import:
```go
import (
    // ... existing imports ...
    "github.com/igorzel/mytets/internal/commands/one"
)
```

In the command setup (find where `version_cmd` is added), add:
```go
root.AddCommand(one.New(cfg))
```

---

## 6. Build and Test Locally

### Build
```bash
go build -o mytets ./cmd/mytets/main.go
```

### Run (plain text)
```bash
./mytets one
```

**Expected output**: A random phrase from the file

### Run (JSON)
```bash
./mytets --output json one
```

**Expected output**: JSON object like `{"message":"..."}`

### Run multiple times
```bash
./mytets one
./mytets one
./mytets one
```

**Expected**: Phrases may differ (demonstrating randomness)

---

## 7. Run Unit Tests

```bash
# Test phrases package
go test ./internal/phrases/...

# Test one command
go test ./internal/commands/one/...
```

**Expected**: All tests pass

---

## 8. Create Integration Tests

**File**: `tests/integration/one_command_test.go`

```go
package integration

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/igorzel/mytets/internal/cli"
)

func TestOneCommandPlainText(t *testing.T) {
	stdout, stderr, exitCode := cli.ExecuteArgs([]string{"one"})

	if exitCode != 0 {
		t.Errorf("exit code = %d, want 0; stderr: %s", exitCode, stderr)
	}

	if stdout == "" {
		t.Error("stdout should not be empty")
	}

	// Verify it's not JSON
	if bytes.HasPrefix([]byte(stdout), []byte("{")) {
		t.Error("plain text output should not be JSON")
	}

	if stderr != "" {
		t.Errorf("stderr should be empty, got %q", stderr)
	}
}

func TestOneCommandJSON(t *testing.T) {
	stdout, stderr, exitCode := cli.ExecuteArgs([]string{"--output", "json", "one"})

	if exitCode != 0 {
		t.Errorf("exit code = %d, want 0; stderr: %s", exitCode, stderr)
	}

	// Verify JSON is valid
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &parsed); err != nil {
		t.Fatalf("stdout is not valid JSON: %v", err)
	}

	if _, ok := parsed["message"]; !ok {
		t.Error("JSON output missing 'message' field")
	}

	if stderr != "" {
		t.Errorf("stderr should be empty, got %q", stderr)
	}
}

func TestOneCommandInvalidFormat(t *testing.T) {
	stdout, stderr, exitCode := cli.ExecuteArgs([]string{"--output", "xml", "one"})

	if exitCode != 1 {
		t.Errorf("exit code = %d, want 1", exitCode)
	}

	if stdout != "" {
		t.Errorf("stdout should be empty on error, got %q", stdout)
	}

	if stderr == "" {
		t.Error("stderr should contain error message")
	}

	if !strings.Contains(strings.ToLower(stderr), "format") {
		t.Errorf("stderr should mention format, got %q", stderr)
	}
}

func TestOneCommandRandomness(t *testing.T) {
	results := make(map[string]int)

	for i := 0; i < 50; i++ {
		stdout, _, exitCode := cli.ExecuteArgs([]string{"one"})
		if exitCode != 0 {
			t.Fatalf("run %d: exit code = %d, want 0", i, exitCode)
		}
		results[strings.TrimSpace(stdout)]++
	}

	if len(results) == 1 {
		t.Errorf("expected multiple different phrases across 50 runs, got only 1")
	}
}
```

### Run integration tests
```bash
go test ./tests/integration/...
```

**Expected**: All tests pass

---

## 9. Run All Tests

```bash
go test ./...
```

**Expected**: All unit and integration tests pass

---

## 10. Verify Complete Build

```bash
go build -o mytets ./cmd/mytets/main.go
./mytets --help
./mytets one --help
./mytets one
./mytets --output json one
```

---

## 11. Next Steps

### Documentation
- [ ] Update README.md with `one` command usage example
- [ ] Add help text to command (already in code via `.Short` and `.Long`)

### Review
- [ ] Code review for `internal/phrases/` package
- [ ] Code review for `internal/commands/one/` command
- [ ] Verify integration tests cover all scenarios

### Performance Testing (Optional)
```bash
# Measure startup + execution time
time ./mytets one
time ./mytets --output json one
```

**Expected**: <100ms total (typically 5-15ms)

---

## Troubleshooting

### "embedded file phrases.json not found"
**Cause**: File path in `//go:embed` doesn't match actual file location  
**Fix**: Verify `embedded/phrases.json` exists relative to `internal/phrases/phrases.go`

### "json: cannot unmarshal X into Y"
**Cause**: JSON file format doesn't match struct definition  
**Fix**: Verify `embedded/phrases.json` has `"messages"` array with `"text"` fields

### "invalid memory address" panic
**Cause**: `phrases` variable is nil (init panic wasn't caught)  
**Fix**: Check application startup logs for panic message

### Tests fail with "flag provided but not defined"
**Cause**: `--output` flag not properly registered in root command  
**Fix**: Verify flag is added in `cli/root.go` or in the one command setup

---

## Key Files Reference

| File | Purpose | Status |
|------|---------|--------|
| `internal/phrases/phrases.json` | Phrase data source | Create |
| `internal/phrases/phrases.go` | Phrase loading & random selection | Create |
| `internal/phrases/phrases_test.go` | Unit tests for phrases | Create |
| `internal/commands/one/one.go` | One command implementation | Create |
| `internal/commands/one/one_test.go` | Unit tests for command | Create |
| `tests/integration/one_command_test.go` | Integration tests | Create |
| `internal/cli/root.go` | Register one command | Modify |
| `README.md` | Project documentation | Update (optional) |

---

## Performance Checklist

- [ ] `go build` completes in <5 seconds
- [ ] `./mytets one` executes in <100ms
- [ ] `./mytets --output json one` executes in <100ms
- [ ] Random phrase selection works across multiple runs
- [ ] All tests pass

---

## Success Criteria

✅ Implementation complete when:
- [ ] `mytets one` returns random phrase (plain text)
- [ ] `mytets --output json one` returns valid, compact JSON
- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] No compiler warnings
- [ ] Performance <100ms on modern hardware
