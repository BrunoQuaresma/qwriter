package cli

import (
	"io"
	"os"

	"github.com/BrunoQuaresma/openwritter/pkg/owriter"
	"github.com/spf13/cobra"
)

type CLI struct {
	writer owriter.Writer
	cmd    *cobra.Command
	stdout io.Writer
	stderr io.Writer
}

type Options struct {
	Writer owriter.Writer
	Stdout io.Writer
	Stderr io.Writer
}

type Flag string

const (
	FlagFiles Flag = "files"
)

func New(o Options) *CLI {
	if o.Writer == nil {
		panic("Writer is required")
	}
	if o.Stdout == nil {
		o.Stdout = os.Stdout
	}
	if o.Stderr == nil {
		o.Stderr = os.Stderr
	}

	// Setup
	cli := &CLI{
		writer: o.Writer,
		stdout: o.Stdout,
		stderr: o.Stderr,
	}
	cli.cmd = cli.root()
	cli.cmd.SetOut(cli.stdout)
	cli.cmd.SetErr(cli.stderr)

	// Add commands
	cli.cmd.AddCommand(cli.review())

	return cli
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
