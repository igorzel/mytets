package phrases

import "testing"

func contains(values []string, target string) bool {
	for _, v := range values {
		if v == target {
			return true
		}
	}
	return false
}

func TestMessagesLoaded(t *testing.T) {
	messages := Messages()
	if len(messages) == 0 {
		t.Fatal("expected embedded messages to be loaded")
	}
	for i, m := range messages {
		if m == "" {
			t.Fatalf("message at index %d is empty", i)
		}
	}
}

func TestRandomMessageFromEmbeddedSet(t *testing.T) {
	messages := Messages()
	if len(messages) == 0 {
		t.Fatal("expected non-empty embedded messages")
	}

	got, err := RandomMessage()
	if err != nil {
		t.Fatalf("RandomMessage returned error: %v", err)
	}
	if !contains(messages, got) {
		t.Fatalf("RandomMessage returned unexpected value %q", got)
	}
}

func TestRandomMessageVariesAcrossRuns(t *testing.T) {
	messages := Messages()
	if len(messages) < 2 {
		t.Skip("embedded data has fewer than two messages")
	}

	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		got, err := RandomMessage()
		if err != nil {
			t.Fatalf("RandomMessage returned error: %v", err)
		}
		seen[got] = true
	}

	if len(seen) < 2 {
		t.Fatalf("expected at least 2 distinct messages in 100 runs, got %d", len(seen))
	}
}
