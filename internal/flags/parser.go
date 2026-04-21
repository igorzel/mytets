// Package flags provides CLI parser configuration types consumed by
// internal/cli when constructing the Cobra command tree.
package flags

import "fmt"

// OutputFormat represents the supported output serialisation formats.
type OutputFormat string

const (
	// OutputFormatText is the default human-readable plain-text format.
	OutputFormatText OutputFormat = "text"
	// OutputFormatJSON is the machine-readable JSON format (--output json).
	OutputFormatJSON OutputFormat = "json"
)

// ParserConfig holds validated configuration for the CLI command tree.
type ParserConfig struct {
	// Output is the requested output format (default: text).
	Output OutputFormat
}

// NewParserConfig returns a ParserConfig with defaults applied.
func NewParserConfig() ParserConfig {
	return ParserConfig{
		Output: OutputFormatText,
	}
}

// ValidOutputFormats lists all format values accepted by --output / -o.
var ValidOutputFormats = []OutputFormat{OutputFormatText, OutputFormatJSON}

// ParseOutputFormat validates and returns an OutputFormat from a raw string.
// An error is returned for unrecognised values.
func ParseOutputFormat(raw string) (OutputFormat, error) {
	switch OutputFormat(raw) {
	case OutputFormatText:
		return OutputFormatText, nil
	case OutputFormatJSON:
		return OutputFormatJSON, nil
	default:
		return "", fmt.Errorf("unsupported output format %q: must be one of text, json", raw)
	}
}
