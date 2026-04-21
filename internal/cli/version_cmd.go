package cli

import (
	"encoding/json"
	"fmt"

	"github.com/igorzel/mytets/internal/flags"
	"github.com/igorzel/mytets/internal/version"
	"github.com/spf13/cobra"
)

// versionOutput is the JSON envelope written by --output json.
type versionOutput struct {
	Version string `json:"version"`
}

// newVersionCmd builds the `version` subcommand. The command accepts an
// optional --output / -o flag; plain text is the default.
func newVersionCmd(cfg flags.ParserConfig) *cobra.Command {
	var outputRaw string

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the application version and exit",
		// No positional arguments are accepted.
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			format, err := flags.ParseOutputFormat(outputRaw)
			if err != nil {
				return err
			}

			switch format {
			case flags.OutputFormatJSON:
				out, jsonErr := json.Marshal(versionOutput{Version: version.Version})
				if jsonErr != nil {
					return fmt.Errorf("failed to encode JSON: %w", jsonErr)
				}
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), string(out))
			default:
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), version.Version)
			}
			return nil
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
