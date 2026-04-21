package phrases

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

//go:embed phrases.json
var phrasesJSON []byte

type phraseEntry struct {
	Text string `json:"text"`
}

type phraseDocument struct {
	Messages []phraseEntry `json:"messages"`
}

var (
	doc     phraseDocument
	loadErr error
	rng     = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func init() {
	loadErr = json.Unmarshal(phrasesJSON, &doc)
	if loadErr != nil {
		return
	}
	if len(doc.Messages) == 0 {
		loadErr = errors.New("embedded phrases contain no messages")
		return
	}
	for i, m := range doc.Messages {
		if strings.TrimSpace(m.Text) == "" {
			loadErr = fmt.Errorf("embedded phrase at index %d is empty", i)
			return
		}
	}
}

// Messages returns a copy of all available phrase texts.
func Messages() []string {
	if loadErr != nil {
		return nil
	}
	out := make([]string, 0, len(doc.Messages))
	for _, m := range doc.Messages {
		out = append(out, m.Text)
	}
	return out
}

// RandomMessage returns one randomly selected message from the embedded set.
func RandomMessage() (string, error) {
	if loadErr != nil {
		return "", fmt.Errorf("failed to load embedded phrases: %w", loadErr)
	}
	if len(doc.Messages) == 0 {
		return "", errors.New("no phrases available")
	}
	return doc.Messages[rng.Intn(len(doc.Messages))].Text, nil
}
