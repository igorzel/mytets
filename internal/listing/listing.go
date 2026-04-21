package listing

import "math/rand"

// Select picks up to count unique phrases from the given slice using the
// provided random source. If count exceeds the number of unique phrases,
// all unique phrases are returned. Duplicate entries in the input are
// deduplicated before selection.
func Select(phrases []string, count int, rng *rand.Rand) []string {
	if count <= 0 || len(phrases) == 0 {
		return nil
	}

	unique := deduplicate(phrases)

	if count > len(unique) {
		count = len(unique)
	}

	// Fisher-Yates partial shuffle: shuffle the first `count` positions.
	for i := 0; i < count; i++ {
		j := i + rng.Intn(len(unique)-i)
		unique[i], unique[j] = unique[j], unique[i]
	}

	return unique[:count]
}

func deduplicate(input []string) []string {
	seen := make(map[string]struct{}, len(input))
	result := make([]string, 0, len(input))
	for _, s := range input {
		if _, exists := seen[s]; !exists {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
}
