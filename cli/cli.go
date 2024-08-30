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
	config *Config
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

	// Flags
	var configPath string

	// Setup
	cli := &CLI{
		writer: o.Writer,
		stdout: o.Stdout,
		stderr: o.Stderr,
		config: nil,
	}
	cli.cmd = &cobra.Command{
		Use:   "qwriter [text]",
		Short: "QWriter CLI is a tool to generate and write text using OpenAI's GPT-4 model.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath != "" {
				config, err := NewConfig(configPath)
				if err != nil {
					return fmt.Errorf("failed to load configuration file: %w", err)
				}
				cli.config = config
			}

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
	cli.cmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to the QWriter CLI configuration file")

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
