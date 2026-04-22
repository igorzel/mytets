package cli

import (
	"fmt"
	"strings"

	"github.com/igorzel/mytets/internal/commands/list"
	"github.com/igorzel/mytets/internal/commands/one"
	"github.com/igorzel/mytets/internal/flags"
	"github.com/igorzel/mytets/internal/i18n"
	"github.com/spf13/cobra"
)

// newRootCmd builds and returns the Cobra root command using the supplied
// parser configuration. All subcommands are registered here.
func newRootCmd(cfg flags.ParserConfig) *cobra.Command {
	root := &cobra.Command{
		Use:   "mytets",
		Short: i18n.Translate("root.short"),
		// Silence default error printing; errors are handled by Execute/ExecuteArgs.
		SilenceErrors: true,
		SilenceUsage:  true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	root.AddCommand(newVersionCmd(cfg))
	root.AddCommand(one.New(cfg))
	root.AddCommand(list.New(cfg))

	root.SetUsageTemplate(localizedUsageTemplate())
	root.SetFlagErrorFunc(flagErrorFunc)
	applyHelpFlags(root)

	return root
}

// localizedUsageTemplate returns a Cobra usage template with all structural
// labels replaced by their translations from the active locale.
func localizedUsageTemplate() string {
	usage := i18n.Translate("help.usage")
	cmds := i18n.Translate("help.available_commands")
	flg := i18n.Translate("help.flags")
	gflg := i18n.Translate("help.global_flags")
	aliases := i18n.Translate("help.aliases")
	cmdWord := i18n.Translate("help.command_word")
	additional := fmt.Sprintf(i18n.Translate("help.additional_help"), "{{.CommandPath}}")

	return usage + `{{if .Runnable}}
  {{localUseLine .}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [` + cmdWord + `]{{end}}{{if gt (len .Aliases) 0}}

` + aliases + `
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

` + cmds + `{{range .Commands}}{{if and .IsAvailableCommand (ne .Name "help")}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

` + flg + `
{{localFlagUsages .LocalFlags | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

` + gflg + `
{{localFlagUsages .InheritedFlags | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableSubCommands}}

` + additional + `{{end}}
`
}

// applyHelpFlags walks the command tree and sets localized help-flag usage
// text on every command (e.g., "help for mytets" → "довідка для mytets").
func applyHelpFlags(cmd *cobra.Command) {
	cmd.InitDefaultHelpFlag()
	if f := cmd.Flags().Lookup("help"); f != nil {
		f.Usage = fmt.Sprintf(i18n.Translate("help.help_for"), cmd.Name())
	}
	for _, child := range cmd.Commands() {
		applyHelpFlags(child)
	}
}

func init() {
	cobra.AddTemplateFuncs(localTemplateFuncs())
}

func localTemplateFuncs() map[string]any {
	return map[string]any{
		"localUseLine": func(cmd *cobra.Command) string {
			line := cmd.UseLine()
			flgWord := i18n.Translate("help.flags_word")
			if flgWord != "flags" {
				line = strings.ReplaceAll(line, "[flags]", "["+flgWord+"]")
			}
			return line
		},
		"localFlagUsages": func(fs interface{ FlagUsages() string }) string {
			usage := fs.FlagUsages()
			defLabel := i18n.Translate("help.default_label")
			if defLabel != "default" {
				usage = strings.ReplaceAll(usage, "(default ", "("+defLabel+" ")
			}
			return usage
		},
	}
}

// flagErrorFunc intercepts flag parsing errors and returns localized messages.
func flagErrorFunc(_ *cobra.Command, err error) error {
	msg := err.Error()
	if strings.HasPrefix(msg, "unknown flag: ") {
		flagName := strings.TrimPrefix(msg, "unknown flag: ")
		return fmt.Errorf(i18n.Translate("error.unknown_flag"), flagName)
	}
	return err
}

// translateExecError translates Cobra execution errors (e.g., unknown command)
// into localized messages. Non-matching errors pass through unchanged.
func translateExecError(err error) error {
	msg := err.Error()
	// Cobra format: `unknown command "X" for "Y"`
	if strings.HasPrefix(msg, `unknown command "`) {
		parts := strings.SplitN(msg, `"`, 5)
		if len(parts) >= 4 {
			cmdName := parts[1]
			parentName := parts[3]
			return fmt.Errorf(i18n.Translate("error.unknown_command"), cmdName, parentName)
		}
	}
	return err
}
