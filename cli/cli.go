package cli

import (
	"errors"
	"fmt"
	"io"

	"github.com/BrunoQuaresma/openwritter/pkg/qwriter"
	"github.com/spf13/cobra"
)

type CLI struct {
	writer qwriter.Writer
	cmd    *cobra.Command
	stdout io.Writer
	stderr io.Writer
}

type Options struct {
	Writer qwriter.Writer
	Stdout io.Writer
	Stderr io.Writer
}

func New(o Options) (*CLI, error) {
	if o.Writer == nil {
		return nil, errors.New("cli.Options.Writer is required")
	}
	if o.Stdout == nil {
		return nil, errors.New("cli.Options.Stdout is required")
	}
	if o.Stderr == nil {
		return nil, errors.New("cli.Options.Stderr is required")
	}

	// Setup
	cli := &CLI{
		writer: o.Writer,
		stdout: o.Stdout,
		stderr: o.Stderr,
	}
	cli.cmd = &cobra.Command{
		Use:   "qwriter [text]",
		Short: "OpenWriter CLI is a tool to generate and write text using OpenAI's GPT-4 model.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txt := args[0]
			s, err := cli.writer.Suggestions(txt)
			if err != nil {
				return fmt.Errorf("failed to get suggestions for %s: %w", txt, err)
			}
			txt, err = cli.writer.Apply(txt, s)
			if err != nil {
				return fmt.Errorf("failed to apply suggestions for %s: %w", txt, err)
			}
			fmt.Fprintln(cli.stdout, txt)
			return nil
		},
	}
	cli.cmd.SetOut(cli.stdout)
	cli.cmd.SetErr(cli.stderr)

	return cli, nil
}

func (cli *CLI) Execute() error {
	return cli.cmd.Execute()
}

// Run executes the CLI with the provided arguments. This function is primarily
// used for testing purposes.
func (cli *CLI) Run(args []string) error {
	cli.cmd.SetArgs(args)
	return cli.Execute()
}
