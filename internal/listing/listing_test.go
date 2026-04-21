package listing

import (
	"math/rand"
	"testing"
)

func TestSelectDefaultCount(t *testing.T) {
	phrases := []string{"a", "b", "c", "d", "e", "f"}
	rng := rand.New(rand.NewSource(42))

	result := Select(phrases, 5, rng)

	if len(result) != 5 {
		t.Fatalf("Select(6 phrases, 5) returned %d items, want 5", len(result))
	}
	assertUnique(t, result)
	assertSubset(t, result, phrases)
}

func TestSelectCustomCount(t *testing.T) {
	phrases := []string{"a", "b", "c", "d", "e"}
	rng := rand.New(rand.NewSource(42))

	result := Select(phrases, 3, rng)

	if len(result) != 3 {
		t.Fatalf("Select(5 phrases, 3) returned %d items, want 3", len(result))
	}
	assertUnique(t, result)
	assertSubset(t, result, phrases)
}

func TestSelectOversizedCountCap(t *testing.T) {
	phrases := []string{"a", "b", "c"}
	rng := rand.New(rand.NewSource(42))

	result := Select(phrases, 10, rng)

	if len(result) != 3 {
		t.Fatalf("Select(3 phrases, 10) returned %d items, want 3", len(result))
	}
	assertUnique(t, result)
	assertSubset(t, result, phrases)
}

func TestSelectExactCount(t *testing.T) {
	phrases := []string{"a", "b", "c"}
	rng := rand.New(rand.NewSource(42))

	result := Select(phrases, 3, rng)

	if len(result) != 3 {
		t.Fatalf("Select(3 phrases, 3) returned %d items, want 3", len(result))
	}
	assertUnique(t, result)
	assertSubset(t, result, phrases)
}

func TestSelectUniquenessGuarantee(t *testing.T) {
	phrases := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	rng := rand.New(rand.NewSource(99))

	for i := 0; i < 50; i++ {
		result := Select(phrases, 5, rng)
		assertUnique(t, result)
	}
}

func TestSelectDeduplicatesDuplicateSourceEntries(t *testing.T) {
	phrases := []string{"a", "b", "a", "c", "b", "d"}
	rng := rand.New(rand.NewSource(42))

	result := Select(phrases, 10, rng)

	if len(result) != 4 {
		t.Fatalf("Select(6 entries, 4 unique, count=10) returned %d items, want 4", len(result))
	}
	assertUnique(t, result)
}

func TestSelectDeterministicSeam(t *testing.T) {
	phrases := []string{"a", "b", "c", "d", "e"}

	result1 := Select(phrases, 3, rand.New(rand.NewSource(42)))
	result2 := Select(phrases, 3, rand.New(rand.NewSource(42)))

	if len(result1) != len(result2) {
		t.Fatalf("deterministic runs returned different lengths: %d vs %d", len(result1), len(result2))
	}
	for i := range result1 {
		if result1[i] != result2[i] {
			t.Errorf("deterministic runs differ at index %d: %q vs %q", i, result1[i], result2[i])
		}
	}
}

func TestSelectCountZero(t *testing.T) {
	phrases := []string{"a", "b", "c"}
	rng := rand.New(rand.NewSource(42))

	result := Select(phrases, 0, rng)

	if result != nil {
		t.Fatalf("Select(count=0) returned %v, want nil", result)
	}
}

func TestSelectEmptyInput(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	result := Select(nil, 5, rng)
	if result != nil {
		t.Fatalf("Select(nil, 5) returned %v, want nil", result)
	}

	result = Select([]string{}, 5, rng)
	if result != nil {
		t.Fatalf("Select(empty, 5) returned %v, want nil", result)
	}
}

func TestSelectNegativeCount(t *testing.T) {
	phrases := []string{"a", "b", "c"}
	rng := rand.New(rand.NewSource(42))

	result := Select(phrases, -1, rng)

	if result != nil {
		t.Fatalf("Select(count=-1) returned %v, want nil", result)
	}
}

func TestSelectSinglePhrase(t *testing.T) {
	phrases := []string{"only"}
	rng := rand.New(rand.NewSource(42))

	result := Select(phrases, 1, rng)

	if len(result) != 1 || result[0] != "only" {
		t.Fatalf("Select(1 phrase, 1) = %v, want [only]", result)
	}
}

// --- helpers ---

func assertUnique(t *testing.T, items []string) {
	t.Helper()
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		if _, exists := seen[item]; exists {
			t.Errorf("duplicate item found: %q", item)
		}
		seen[item] = struct{}{}
	}
}

func assertSubset(t *testing.T, items, superset []string) {
	t.Helper()
	allowed := make(map[string]struct{}, len(superset))
	for _, s := range superset {
		allowed[s] = struct{}{}
	}
	for _, item := range items {
		if _, ok := allowed[item]; !ok {
			t.Errorf("item %q not found in allowed set", item)
		}
	}
}
