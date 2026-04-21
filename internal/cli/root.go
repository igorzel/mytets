package cli

import (
	"github.com/igorzel/mytets/internal/commands/list"
	"github.com/igorzel/mytets/internal/commands/one"
	"github.com/igorzel/mytets/internal/flags"
	"github.com/spf13/cobra"
)

// newRootCmd builds and returns the Cobra root command using the supplied
// parser configuration. All subcommands are registered here.
func newRootCmd(cfg flags.ParserConfig) *cobra.Command {
	root := &cobra.Command{
		Use:   "mytets",
		Short: "mytets — a lightweight CLI tool",
		// Silence default error printing; errors are handled by Execute/ExecuteArgs.
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	root.AddCommand(newVersionCmd(cfg))
	root.AddCommand(one.New(cfg))
	root.AddCommand(list.New(cfg))

	return root
}
