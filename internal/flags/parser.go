// Package flags provides CLI parser configuration types consumed by
// internal/cli when constructing the Cobra command tree.
package flags

import (
	"fmt"
	"strings"
)

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

// ConsumeLeadingGlobalOutput parses and removes leading global output flags
// from CLI args (for example: --output json one). The returned slice preserves
// remaining args in original order.
func ConsumeLeadingGlobalOutput(args []string) (OutputFormat, []string, error) {
	format := OutputFormatText
	i := 0

	for i < len(args) {
		arg := args[i]
		switch {
		case arg == "--output" || arg == "-o":
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("missing value for %s", arg)
			}
			parsed, err := ParseOutputFormat(args[i+1])
			if err != nil {
				return "", nil, err
			}
			format = parsed
			i += 2
		case strings.HasPrefix(arg, "--output="):
			parsed, err := ParseOutputFormat(strings.TrimPrefix(arg, "--output="))
			if err != nil {
				return "", nil, err
			}
			format = parsed
			i++
		case strings.HasPrefix(arg, "-o="):
			parsed, err := ParseOutputFormat(strings.TrimPrefix(arg, "-o="))
			if err != nil {
				return "", nil, err
			}
			format = parsed
			i++
		default:
			return format, args[i:], nil
		}
	}

	return format, []string{}, nil
}
