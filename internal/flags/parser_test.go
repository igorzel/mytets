package flags_test

import (
	"testing"

	"github.com/igorzel/mytets/internal/flags"
	"github.com/igorzel/mytets/internal/i18n"
)

func init() {
	i18n.LoadBundle()
	i18n.SetLang("en")
}

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

func TestConsumeLeadingGlobalOutput(t *testing.T) {
	tests := []struct {
		name       string
		in         []string
		wantFormat flags.OutputFormat
		wantArgs   []string
		wantErr    bool
	}{
		{
			name:       "no-global-flag",
			in:         []string{"one"},
			wantFormat: flags.OutputFormatText,
			wantArgs:   []string{"one"},
		},
		{
			name:       "global-long-form",
			in:         []string{"--output", "json", "one"},
			wantFormat: flags.OutputFormatJSON,
			wantArgs:   []string{"one"},
		},
		{
			name:       "global-short-equals",
			in:         []string{"-o=json", "one"},
			wantFormat: flags.OutputFormatJSON,
			wantArgs:   []string{"one"},
		},
		{
			name:    "missing-value",
			in:      []string{"--output"},
			wantErr: true,
		},
		{
			name:    "bad-value",
			in:      []string{"--output", "yaml", "one"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotFormat, gotArgs, err := flags.ConsumeLeadingGlobalOutput(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotFormat != tc.wantFormat {
				t.Fatalf("format = %q, want %q", gotFormat, tc.wantFormat)
			}
			if len(gotArgs) != len(tc.wantArgs) {
				t.Fatalf("args len = %d, want %d", len(gotArgs), len(tc.wantArgs))
			}
			for i := range gotArgs {
				if gotArgs[i] != tc.wantArgs[i] {
					t.Fatalf("args[%d] = %q, want %q", i, gotArgs[i], tc.wantArgs[i])
				}
			}
		})
	}
}

// T016: Verify Ukrainian unsupported format error.
func TestParseOutputFormatUkrainian(t *testing.T) {
	i18n.SetLang("uk")
	defer i18n.SetLang("en")

	_, err := flags.ParseOutputFormat("xml")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if want := `непідтримуваний формат виводу: "xml"`; err.Error() != want {
		t.Errorf("error = %q, want %q", err.Error(), want)
	}
}
