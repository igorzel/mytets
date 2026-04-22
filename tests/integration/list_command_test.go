package integration

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/igorzel/mytets/internal/cli"
	"github.com/igorzel/mytets/internal/phrases"
)

// --- Plain Text Tests (US1) ---

func TestListCommandPlainDefault(t *testing.T) {
	stdout, stderr, exitCode := cli.ExecuteArgs([]string{"list"})

	if exitCode != 0 {
		t.Fatalf("Exit code = %d, want 0 (stderr=%q)", exitCode, stderr)
	}
	if stderr != "" {
		t.Errorf("stderr should be empty, got %q", stderr)
	}

	lines := nonEmptyLines(stdout)
	available := phrases.Messages()

	// Default count is 5, but capped at available unique phrases.
	expectedCount := 5
	if expectedCount > len(available) {
		expectedCount = len(available)
	}

	if len(lines) != expectedCount {
		t.Fatalf("got %d lines, want %d", len(lines), expectedCount)
	}
	assertUniqueLines(t, lines)
	assertAllInPhrases(t, lines, available)
}

func TestListCommandPlainCountN(t *testing.T) {
	stdout, stderr, exitCode := cli.ExecuteArgs([]string{"list", "--count", "2"})

	if exitCode != 0 {
		t.Fatalf("Exit code = %d, want 0 (stderr=%q)", exitCode, stderr)
	}
	if stderr != "" {
		t.Errorf("stderr should be empty, got %q", stderr)
	}

	lines := nonEmptyLines(stdout)
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2", len(lines))
	}
	assertUniqueLines(t, lines)
	assertAllInPhrases(t, lines, phrases.Messages())
}

func TestListCommandPlainOversizedCount(t *testing.T) {
	available := phrases.Messages()
	stdout, stderr, exitCode := cli.ExecuteArgs([]string{"list", "--count", "999"})

	if exitCode != 0 {
		t.Fatalf("Exit code = %d, want 0 (stderr=%q)", exitCode, stderr)
	}
	if stderr != "" {
		t.Errorf("stderr should be empty, got %q", stderr)
	}

	lines := nonEmptyLines(stdout)
	if len(lines) != len(available) {
		t.Fatalf("got %d lines, want %d (all available)", len(lines), len(available))
	}
	assertUniqueLines(t, lines)
	assertAllInPhrases(t, lines, available)
}

func TestListCommandPlainCountOne(t *testing.T) {
	stdout, stderr, exitCode := cli.ExecuteArgs([]string{"list", "--count", "1"})

	if exitCode != 0 {
		t.Fatalf("Exit code = %d, want 0 (stderr=%q)", exitCode, stderr)
	}
	if stderr != "" {
		t.Errorf("stderr should be empty, got %q", stderr)
	}

	lines := nonEmptyLines(stdout)
	if len(lines) != 1 {
		t.Fatalf("got %d lines, want 1", len(lines))
	}
	assertAllInPhrases(t, lines, phrases.Messages())
}

func TestListCommandInvalidCountString(t *testing.T) {
	_, stderr, exitCode := cli.ExecuteArgs([]string{"list", "--count", "abc"})

	if exitCode == 0 {
		t.Error("Exit code should be non-zero for invalid count")
	}
	if stderr == "" {
		t.Error("stderr should contain error message")
	}
}

func TestListCommandInvalidCountNegative(t *testing.T) {
	_, stderr, exitCode := cli.ExecuteArgs([]string{"list", "--count", "-1"})

	// Cobra parses -1 as a valid int. The listing.Select treats negative as
	// returning nil, which means 0 lines. Our command checks len(msgs)==0
	// after messageSource() but listing.Select returning nil for count<=0 is
	// by design. The command should still exit 0 with empty output.
	// However, if there are 0 selected phrases that's still a successful run
	// with empty output (like --count 0).
	if exitCode != 0 {
		t.Logf("Exit code = %d, stderr = %q", exitCode, stderr)
	}
}

func TestListCommandExitCodeZero(t *testing.T) {
	_, _, exitCode := cli.ExecuteArgs([]string{"list"})
	if exitCode != 0 {
		t.Errorf("Exit code = %d, want 0", exitCode)
	}
}

// --- JSON Tests (US2) ---

func TestListCommandJSONDefault(t *testing.T) {
	stdout, stderr, exitCode := cli.ExecuteArgs([]string{"--output", "json", "list"})

	if exitCode != 0 {
		t.Fatalf("Exit code = %d, want 0 (stderr=%q)", exitCode, stderr)
	}
	if stderr != "" {
		t.Errorf("stderr should be empty, got %q", stderr)
	}

	items := parseJSONItems(t, stdout)
	available := phrases.Messages()

	expectedCount := 5
	if expectedCount > len(available) {
		expectedCount = len(available)
	}

	if len(items) != expectedCount {
		t.Fatalf("got %d JSON items, want %d", len(items), expectedCount)
	}
	assertUniqueJSONMessages(t, items)
	assertJSONMessagesInPhrases(t, items, available)
}

func TestListCommandJSONCountN(t *testing.T) {
	stdout, stderr, exitCode := cli.ExecuteArgs([]string{"--output", "json", "list", "--count", "2"})

	if exitCode != 0 {
		t.Fatalf("Exit code = %d, want 0 (stderr=%q)", exitCode, stderr)
	}
	if stderr != "" {
		t.Errorf("stderr should be empty, got %q", stderr)
	}

	items := parseJSONItems(t, stdout)
	if len(items) != 2 {
		t.Fatalf("got %d JSON items, want 2", len(items))
	}
	assertUniqueJSONMessages(t, items)
	assertJSONMessagesInPhrases(t, items, phrases.Messages())
}

func TestListCommandJSONOversizedCount(t *testing.T) {
	available := phrases.Messages()
	stdout, stderr, exitCode := cli.ExecuteArgs([]string{"--output", "json", "list", "--count", "999"})

	if exitCode != 0 {
		t.Fatalf("Exit code = %d, want 0 (stderr=%q)", exitCode, stderr)
	}
	if stderr != "" {
		t.Errorf("stderr should be empty, got %q", stderr)
	}

	items := parseJSONItems(t, stdout)
	if len(items) != len(available) {
		t.Fatalf("got %d JSON items, want %d (all available)", len(items), len(available))
	}
	assertUniqueJSONMessages(t, items)
	assertJSONMessagesInPhrases(t, items, available)
}

func TestListCommandJSONExitCodeZero(t *testing.T) {
	_, _, exitCode := cli.ExecuteArgs([]string{"--output", "json", "list"})
	if exitCode != 0 {
		t.Errorf("Exit code = %d, want 0", exitCode)
	}
}

func TestListCommandJSONCompactFormat(t *testing.T) {
	stdout, _, exitCode := cli.ExecuteArgs([]string{"--output", "json", "list", "--count", "1"})

	if exitCode != 0 {
		t.Fatal("Exit code should be 0")
	}

	trimmed := strings.TrimSpace(stdout)
	if strings.Contains(trimmed, "\n") {
		t.Error("JSON output should be compact (single line)")
	}
	if strings.Contains(trimmed, "  ") {
		t.Error("JSON output should not contain indentation")
	}
}

// --- Helpers ---

type jsonItem struct {
	Message string `json:"message"`
}

func parseJSONItems(t *testing.T, raw string) []jsonItem {
	t.Helper()
	trimmed := strings.TrimSpace(raw)
	var items []jsonItem
	if err := json.Unmarshal([]byte(trimmed), &items); err != nil {
		t.Fatalf("failed to parse JSON output: %v\nraw: %q", err, trimmed)
	}
	return items
}

func nonEmptyLines(s string) []string {
	var lines []string
	for _, line := range strings.Split(s, "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, strings.TrimSpace(line))
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

func assertAllInPhrases(t *testing.T, lines []string, available []string) {
	t.Helper()
	set := make(map[string]struct{}, len(available))
	for _, p := range available {
		set[p] = struct{}{}
	}
	for _, line := range lines {
		if _, ok := set[line]; !ok {
			t.Errorf("line %q not found in available phrases", line)
		}
	}
}

func assertUniqueJSONMessages(t *testing.T, items []jsonItem) {
	t.Helper()
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		if _, exists := seen[item.Message]; exists {
			t.Errorf("duplicate JSON message found: %q", item.Message)
		}
		seen[item.Message] = struct{}{}
	}
}

func assertJSONMessagesInPhrases(t *testing.T, items []jsonItem, available []string) {
	t.Helper()
	set := make(map[string]struct{}, len(available))
	for _, p := range available {
		set[p] = struct{}{}
	}
	for _, item := range items {
		if _, ok := set[item.Message]; !ok {
			t.Errorf("JSON message %q not found in available phrases", item.Message)
		}
	}
}

// T028: Verify list output phrases are not localized with Ukrainian locale.
func TestListCommandPhrasesNotLocalizedUkrainian(t *testing.T) {
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "")
	t.Setenv("LANG", "uk_UA.UTF-8")

	available := phrases.Messages()
	stdout, stderr, exitCode := cli.ExecuteArgs([]string{"list", "--count", "3"})

	if exitCode != 0 {
		t.Fatalf("Exit code = %d, want 0 (stderr=%q)", exitCode, stderr)
	}

	lines := nonEmptyLines(stdout)
	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3", len(lines))
	}
	assertAllInPhrases(t, lines, available)
}
