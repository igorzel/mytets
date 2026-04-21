package list

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/igorzel/mytets/internal/flags"
	"github.com/igorzel/mytets/internal/listing"
	"github.com/igorzel/mytets/internal/phrases"
	"github.com/spf13/cobra"
)

var messageSource = phrases.Messages

// New returns a Cobra command for the "list" subcommand.
func New(cfg flags.ParserConfig) *cobra.Command {
	var count int
	var outputRaw string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Display a list of random phrases",
		Long:  "The list command outputs multiple unique random phrases in plain text or JSON format.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			format, err := flags.ParseOutputFormat(outputRaw)
			if err != nil {
				return err
			}

			msgs := messageSource()
			if len(msgs) == 0 {
				return fmt.Errorf("no phrases available")
			}

			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			selected := listing.Select(msgs, count, rng)

			switch format {
			case flags.OutputFormatJSON:
				return outputJSON(cmd, selected)
			default:
				return outputPlain(cmd, selected)
			}
		},
	}

	cmd.Flags().IntVar(&count, "count", 5, "Number of phrases to return")
	cmd.Flags().StringVarP(
		&outputRaw,
		"output", "o",
		string(cfg.Output),
		`Output format: "text" (default) or "json"`,
	)

	return cmd
}

func outputPlain(cmd *cobra.Command, selected []string) error {
	for _, s := range selected {
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), s); err != nil {
			return err
		}
	}
	return nil
}

// Response represents a single phrase in the JSON output array.
type Response struct {
	Message string `json:"message"`
}

func outputJSON(cmd *cobra.Command, selected []string) error {
	items := make([]Response, len(selected))
	for i, s := range selected {
		items[i] = Response{Message: s}
	}
	data, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	_, err = fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return err
}
