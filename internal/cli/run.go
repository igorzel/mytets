package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/igorzel/mytets/internal/flags"
)

// Execute runs the CLI with os.Args[1:] and writes output to os.Stdout /
// os.Stderr. It returns the exit code that the caller should pass to os.Exit.
func Execute() int {
	args := os.Args[1:]
	cfg := flags.NewParserConfig()
	format, rest, err := flags.ConsumeLeadingGlobalOutput(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	cfg.Output = format
	root := newRootCmd(cfg)
	root.SetArgs(rest)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

// ExecuteArgs is a test seam that runs the CLI with the given argument slice
// and captures stdout, stderr, and the exit code without touching os.Stdout,
// os.Stderr, or os.Exit.
func ExecuteArgs(args []string) (stdout, stderr string, exitCode int) {
	cfg := flags.NewParserConfig()
	format, rest, err := flags.ConsumeLeadingGlobalOutput(args)
	if err != nil {
		return "", err.Error() + "\n", 1
	}
	cfg.Output = format
	root := newRootCmd(cfg)

	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	root.SetOut(outBuf)
	root.SetErr(errBuf)
	root.SetArgs(rest)

	if err := root.Execute(); err != nil {
		_, _ = io.WriteString(errBuf, err.Error()+"\n")
		exitCode = 1
	}

	return outBuf.String(), errBuf.String(), exitCode
}
