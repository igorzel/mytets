# Research: One Command - Random Phrase Feature

**Feature**: Implement `mytets one` command to output random phrases from embedded file  
**Date**: 2026-04-21  
**Scope**: Go embedding patterns, JSON unmarshaling, random selection, CLI testing  

## 1. Go embed Package Best Practices

**Decision**: Use `//go:embed` directive with package-level variable  
**Rationale**: 
- Compile-time embedding with zero runtime I/O cost
- Package-level embedded content is immutable and thread-safe
- No external file dependencies—binary is completely self-contained
- Meets <100ms performance requirement
- Standard Go 1.16+ language feature (native, maintained by Go team)

**Implementation Pattern**:
```go
//go:embed phrases.json
var phrasesJSON string  // or []byte for binary data
```

**Error Handling at Initialization**:
```go
var phrases PhraseData

func init() {
    if err := json.Unmarshal([]byte(phrasesJSON), &phrases); err != nil {
        panic(fmt.Sprintf("failed to parse embedded phrases.json: %v", err))
    }
    if len(phrases.Messages) == 0 {
        panic("phrases.json contains no messages")
    }
}
```

**Alternative Rejected**: 
- Pre-1.16 tools like `go-bindata` — now obsolete; native `embed` is superior
- Runtime file loading — violates <100ms requirement and adds deployment complexity

## 2. JSON Unmarshaling Strategy

**Decision**: Unmarshal once at initialization, cache structured data  
**Rationale**:
- CLI tools are stateless and short-lived (single process invocation)
- Parse-once-at-startup eliminates per-command overhead
- Using typed structs with JSON tags is idiomatic Go
- Matches pattern already used in project (version constants)
- Cleaner error handling at startup vs. per-command

**Recommended Package Structure**:
```go
// internal/phrases/phrases.go
package phrases

import (
    "encoding/json"
    "fmt"
)

//go:embed phrases.json
var phrasesJSON string

type data struct {
    Messages []struct {
        Text string `json:"text"`
    } `json:"messages"`
}

var phrases data

func init() {
    if err := json.Unmarshal([]byte(phrasesJSON), &phrases); err != nil {
        panic(fmt.Sprintf("failed to parse embedded phrases.json: %v", err))
    }
    if len(phrases.Messages) == 0 {
        panic("phrases.json contains no messages")
    }
}

// Public API
func GetMessages() []string {
    result := make([]string, len(phrases.Messages))
    for i, m := range phrases.Messages {
        result[i] = m.Text
    }
    return result
}
```

**Alternative Rejected**: 
- Lazy unmarshaling (on first command execution) — unnecessary complexity for short-lived CLI
- Generic `json.RawMessage` — over-engineered; typed structs are idiomatic

## 3. Random Selection Implementation

**Decision**: Use `rand.New(rand.NewSource(...))` with `Intn(len(phrases))`  
**Rationale**:
- Provides uniform distribution across all phrases (O(1) selection)
- Go 1.22+ auto-seeds global random source, but explicit seed is safer
- Your project targets Go 1.25+ compatibility—use explicit seeding for clarity
- Performance cost negligible (<1μs per call)
- Meets <100ms execution target with room to spare

**Go 1.25+ Implementation**:
```go
// internal/phrases/random.go (optional separate file)
package phrases

import (
    "math/rand"
    "time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetRandomPhrase() string {
    messages := GetMessages()
    if len(messages) == 0 {
        return "" // Caught by init() validation
    }
    return messages[rng.Intn(len(messages))]
}
```

**Or inline in phrases.go**: Depends on code organization preference.

**Alternative Rejected**: 
- `crypto/rand` — overkill (10-100x slower); intended for cryptographic randomness, not CLI utilities
- Per-call seeding — inefficient and unnecessary; seed once at startup

## 4. Cobra Command Framework Integration

**Decision**: Use `RunE` with switch on parsed output format, delegate to handler functions  
**Rationale**:
- Matches existing pattern in `version_cmd.go` and `run.go`
- `RunE` returns error automatically sets exit code 1
- Separating output handlers (`outputPlain`, `outputJSON`) keeps code readable
- Using `cmd.OutOrStdout()` is Cobra-idiomatic and testable
- Flag parsing via `flags.ParseOutputFormat()` validates early

**Command Structure** (follows existing pattern):
```go
// internal/commands/one/one.go
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
                return err
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
    fmt.Fprintln(cmd.OutOrStdout()) // Add newline
    return nil
}
```

**Cobra Best Practices Applied**:
- ✅ `Args: cobra.NoArgs` — reject positional arguments
- ✅ `RunE` — return errors; Cobra handles exit codes
- ✅ `cmd.OutOrStdout()` — testable output redirection
- ✅ No custom `Run` field — use `RunE` consistently
- ✅ Avoid `PostRunE` — not needed for simple command

## 5. CLI Integration Testing Strategy

**Decision**: Extend existing `ExecuteArgs` pattern with table-driven subtests  
**Rationale**:
- Project already has `cli.ExecuteArgs()` test seam (perfect!)
- Returns stdout, stderr, and exit code without OS mocking
- Table-driven tests provide readable, maintainable coverage
- Error cases must verify BOTH exit code 1 AND stderr message
- Integration tests belong in `tests/integration/` (existing pattern)

**Test Coverage Checklist**:
- ✅ Exit code 0 on success, 1 on error
- ✅ Stdout contains single phrase (plain text)
- ✅ Stderr empty on success, populated with error message on failure
- ✅ JSON output is valid and compact (single line)
- ✅ JSON has `{"message":"..."}` structure
- ✅ Plain text output is NOT JSON-formatted
- ✅ Invalid output format rejected with error message
- ✅ Randomness across multiple runs (verify >1 unique phrase in 100 runs)
- ✅ Repeated execution produces consistent exit code 0

**Test File Structure**:
```
tests/
└── integration/
    ├── one_command_test.go          # [NEW] One command tests
    └── version_command_test.go      # [EXISTING] Version command
```

**Alternative Rejected**: 
- Shell script integration tests — less maintainable, slower
- Mock-heavy unit tests with testify/mock — your `ExecuteArgs` seam eliminates need for mocks
- No testing — feature completeness requires integration tests per spec

## Summary: Technology Decisions

| Aspect | Pattern | Confidence | Performance Impact |
|--------|---------|------------|-------------------|
| Embedding | `//go:embed` + init() parse | ✅ High | <1ms |
| JSON unmarshaling | Parse once, cache | ✅ High | <1ms |
| Random selection | `math/rand.Intn()` | ✅ High | <1μs |
| Command structure | Cobra `RunE` + handlers | ✅ High | ~1-5ms |
| Output handling | `cmd.OutOrStdout()` with switch | ✅ High | <1ms |
| Testing strategy | `ExecuteArgs` + table-driven | ✅ High | N/A (test code) |

**Total Expected Execution Time**: 5-10ms (well under 100ms target)

---

## Implementation Ready

All key technology decisions have been researched and validated. Patterns align with:
- Go standard library best practices
- Project's existing conventions (version_cmd.go, run.go)
- Constitution principles (performance, reliability, simplicity, extensibility)
- Performance targets (<100ms)

Next phase: Design data model and CLI contracts.
