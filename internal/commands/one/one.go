package one

import (
	"encoding/json"
	"fmt"

	"github.com/igorzel/mytets/internal/flags"
	"github.com/igorzel/mytets/internal/i18n"
	"github.com/igorzel/mytets/internal/phrases"
	"github.com/spf13/cobra"
)

var randomMessage = phrases.RandomMessage

// Response represents the JSON output format for the one command
type Response struct {
	Message string `json:"message"`
}

// New returns a Cobra command for the "one" subcommand
func New(cfg flags.ParserConfig) *cobra.Command {
	var outputRaw string

	cmd := &cobra.Command{
		Use:   "one",
		Short: i18n.Translate("one.short"),
		Long:  i18n.Translate("one.long"),
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
		i18n.Translate("flag.output"),
	)

	return cmd
}

// outputPlain writes the message as plain text to the command output
func outputPlain(cmd *cobra.Command) error {
	msg, err := randomMessage()
	if err != nil {
		return fmt.Errorf(i18n.Translate("error.failed_select_phrase"), err)
	}
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), msg)
	return nil
}

// outputJSON writes the message as compact JSON to the command output
func outputJSON(cmd *cobra.Command) error {
	msg, err := randomMessage()
	if err != nil {
		return fmt.Errorf(i18n.Translate("error.failed_select_phrase"), err)
	}
	resp := Response{Message: msg}
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf(i18n.Translate("error.failed_encode_json"), err)
	}
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}
