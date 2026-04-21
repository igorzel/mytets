package one

import (
	"encoding/json"
	"fmt"

	"github.com/igorzel/mytets/internal/flags"
	"github.com/spf13/cobra"
)

const (
	message = "Fake message, tbd"
)

// Response represents the JSON output format for the one command
type Response struct {
	Message string `json:"message"`
}

// New returns a Cobra command for the "one" subcommand
func New(cfg flags.ParserConfig) *cobra.Command {
	var outputRaw string

	cmd := &cobra.Command{
		Use:   "one",
		Short: "Display the one command message",
		Long:  "The one command outputs a fixed message in plain text or JSON format.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			format, err := flags.ParseOutputFormat(outputRaw)
			if err != nil {
				return err
			}

			switch format {
			case flags.OutputFormatJSON:
				return outputJSON(cmd)
			default:
				return outputPlain(cmd)
			}
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

// outputPlain writes the message as plain text to the command output
func outputPlain(cmd *cobra.Command) error {
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), message)
	return nil
}

// outputJSON writes the message as compact JSON to the command output
func outputJSON(cmd *cobra.Command) error {
	resp := Response{Message: message}
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}
