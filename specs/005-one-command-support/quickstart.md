# Quick Start: Implementing the One Command

This guide walks you through implementing the `one` command for the mytets CLI tool.

## Overview

The `one` command is a simple subcommand that outputs a fixed message, with support for JSON output via the global `--json` flag. The implementation follows existing Cobra patterns and respects the project's Go best practices.

## Prerequisites

- Go 1.26.2 or later
- Existing mytets project with Cobra CLI framework
- Familiarity with Go's `testing` package

## High-Level Implementation Steps

1. **Create the command package** at `internal/commands/one/one.go`
2. **Implement the Cobra command** with plain text and JSON output logic
3. **Register the command** in `internal/cli/root.go`
4. **Write unit tests** for output formatting
5. **Add integration tests** for end-to-end CLI behavior

## Package Structure

```
internal/
├── commands/              # Subcommand implementations (new directory)
│   └── one/               # One command package
│       ├── one.go         # Command implementation
│       └── one_test.go    # Unit tests
└── ... (existing packages)

tests/
└── integration/
    └── one_command_test.go  # Integration tests (new file)
```

## Implementation Outline

### Step 1: Create `internal/commands/one/one.go`

```go
package one

import (
    "encoding/json"
    "fmt"

    "github.com/igorzel/mytets/internal/flags"
    "github.com/spf13/cobra"
)

const (
    message = "Fake message, tbd"
)

type Response struct {
    Message string `json:"message"`
}

// New returns a Cobra command for the "one" subcommand.
func New(cfg flags.ParserConfig) *cobra.Command {
    return &cobra.Command{
        Use:   "one",
        Short: "Display the one command message",
        Long:  "The one command outputs a fixed message in plain text or JSON format.",
        RunE: func(cmd *cobra.Command, args []string) error {
            return execute(cfg)
        },
    }
}

func execute(cfg flags.ParserConfig) error {
    if cfg.OutputJSON {
        return outputJSON()
    }
    return outputPlain()
}

func outputPlain() error {
    fmt.Println(message)
    return nil
}

func outputJSON() error {
    resp := Response{Message: message}
    data, err := json.Marshal(resp)
    if err != nil {
        return fmt.Errorf("failed to marshal JSON: %w", err)
    }
    fmt.Println(string(data))
    return nil
}
```

### Step 2: Create `internal/commands/one/one_test.go`

```go
package one

import (
    "bytes"
    "encoding/json"
    "testing"

    "github.com/igorzel/mytets/internal/flags"
)

func TestExecutePlain(t *testing.T) {
    // Test plain text output
    if err := outputPlain(); err != nil {
        t.Fatalf("outputPlain() error = %v", err)
    }
    // Output assertion would require capturing stdout; see integration tests
}

func TestExecuteJSON(t *testing.T) {
    // Test JSON output format
    if err := outputJSON(); err != nil {
        t.Fatalf("outputJSON() error = %v", err)
    }
}

func TestResponseJSONFormat(t *testing.T) {
    resp := Response{Message: message}
    data, err := json.Marshal(resp)
    if err != nil {
        t.Fatalf("json.Marshal() error = %v", err)
    }

    // Verify JSON is compact and contains expected field
    jsonStr := string(data)
    if jsonStr != `{"message":"Fake message, tbd"}` {
        t.Errorf("JSON format mismatch: got %q, want %q", jsonStr, `{"message":"Fake message, tbd"}`)
    }
}

func TestNew(t *testing.T) {
    cfg := flags.ParserConfig{}
    cmd := New(cfg)
    if cmd == nil {
        t.Fatal("New() returned nil command")
    }
    if cmd.Use != "one" {
        t.Errorf("Command Use = %q, want %q", cmd.Use, "one")
    }
}
```

### Step 3: Register in `internal/cli/root.go`

Modify the `newRootCmd` function to add the one command:

```go
import (
    "github.com/igorzel/mytets/internal/commands/one"  // Add this import
    "github.com/igorzel/mytets/internal/flags"
    "github.com/spf13/cobra"
)

func newRootCmd(cfg flags.ParserConfig) *cobra.Command {
    root := &cobra.Command{
        Use:   "mytets",
        Short: "mytets — a lightweight CLI tool",
        SilenceErrors: true,
        SilenceUsage:  true,
    }

    root.AddCommand(newVersionCmd(cfg))
    root.AddCommand(one.New(cfg))  // Add this line

    return root
}
```

### Step 4: Integration Tests

Create `tests/integration/one_command_test.go`:

```go
package integration

import (
    "bytes"
    "encoding/json"
    "os/exec"
    "strings"
    "testing"
)

func TestOneCommandPlain(t *testing.T) {
    cmd := exec.Command("./main", "one")
    var stdout bytes.Buffer
    cmd.Stdout = &stdout

    if err := cmd.Run(); err != nil {
        t.Fatalf("Command failed: %v", err)
    }

    output := strings.TrimSpace(stdout.String())
    if output != "Fake message, tbd" {
        t.Errorf("Output mismatch: got %q, want %q", output, "Fake message, tbd")
    }
}

func TestOneCommandJSON(t *testing.T) {
    cmd := exec.Command("./main", "--json", "one")
    var stdout bytes.Buffer
    cmd.Stdout = &stdout

    if err := cmd.Run(); err != nil {
        t.Fatalf("Command failed: %v", err)
    }

    var resp map[string]string
    if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
        t.Fatalf("Failed to parse JSON: %v", err)
    }

    if resp["message"] != "Fake message, tbd" {
        t.Errorf("Message mismatch: got %q, want %q", resp["message"], "Fake message, tbd")
    }
}
```

## Key Design Points

1. **Separation of Concerns**: The `one` package is isolated from the root CLI logic, making it easy to test and reuse
2. **Respecting Global Flags**: The command receives `flags.ParserConfig` which includes JSON mode state
3. **No Subcommand-Specific Flags**: The command accepts no special flags; only the global `--json` flag matters
4. **Deterministic Output**: Message and format are hardcoded and identical across runs
5. **Go Best Practices**: Single-letter or short variable names, no stuttering, clean error handling

## Testing Checklist

- [ ] Unit tests verify plain text output function
- [ ] Unit tests verify JSON output function
- [ ] Unit tests verify JSON format is compact and valid
- [ ] Integration tests verify `mytets one` produces correct stdout
- [ ] Integration tests verify `mytets --json one` produces valid JSON with correct message
- [ ] Integration tests verify exit code is 0 for both modes
- [ ] Integration tests verify stderr is empty on success

## Common Pitfalls to Avoid

1. **Pretty-printed JSON**: Ensure `json.Marshal()` is used, not `json.MarshalIndent()`
2. **Missing Trailing Newline**: Use `fmt.Println()`, not `fmt.Print()`, for consistent output
3. **Output to Wrong Stream**: Only write to stdout; stderr should remain empty on success
4. **Flag Handling**: Don't implement standalone `--json` at the subcommand level; rely on global flag
