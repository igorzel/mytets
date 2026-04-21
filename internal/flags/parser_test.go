package flags_test

import (
	"testing"

	"github.com/igorzel/mytets/internal/flags"
)

func TestNewParserConfigDefaults(t *testing.T) {
	cfg := flags.NewParserConfig()
	if cfg.Output != flags.OutputFormatText {
		t.Errorf("expected default output %q, got %q", flags.OutputFormatText, cfg.Output)
	}
}

func TestParseOutputFormat(t *testing.T) {
	tests := []struct {
		input   string
		want    flags.OutputFormat
		wantErr bool
	}{
		{"text", flags.OutputFormatText, false},
		{"json", flags.OutputFormatJSON, false},
		{"", "", true},
		{"yaml", "", true},
		{"XML", "", true},
	}

	for _, tc := range tests {
		got, err := flags.ParseOutputFormat(tc.input)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParseOutputFormat(%q): expected error, got nil", tc.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseOutputFormat(%q): unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.want {
			t.Errorf("ParseOutputFormat(%q): got %q, want %q", tc.input, got, tc.want)
		}
	}
}
