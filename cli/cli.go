package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/BrunoQuaresma/openwritter/pkg/owriter"
	"github.com/bmatcuk/doublestar/v4"
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

	cli := &CLI{
		writer: o.Writer,
		stdout: o.Stdout,
		stderr: o.Stderr,
	}
	cli.cmd = &cobra.Command{
		Use:   "owriter",
		Short: "OpenWriter CLI is a tool to generate and write text using OpenAI's GPT-4 model.",
		Run: func(cmd *cobra.Command, args []string) {
			files, _ := cmd.Flags().GetString(string(FlagFiles))
			if files != "" {
				cli.reviewFiles(files)
				return
			}
		},
	}
	cli.cmd.SetOut(cli.stdout)
	cli.cmd.SetErr(cli.stderr)
	cli.cmd.PersistentFlags().StringP("files", "f", "", "Glob pattern for files to process.")

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

func (cli *CLI) reviewFiles(pattern string) {
	// The CLI only handle files and not directories so we use the WithFilesOnly
	// option to raise an error if a directory is passed.
	matches, err := doublestar.FilepathGlob(pattern, doublestar.WithFilesOnly())
	if err != nil {
		fmt.Fprintln(cli.stderr, pattern, "is not a valid pattern")
		os.Exit(1)
	}
	if matches == nil {
		fmt.Fprintln(cli.stderr, "No files found for pattern", pattern)
		os.Exit(1)
	}
	for _, path := range matches {
		s, err := cli.writer.Suggestions(path)
		if err != nil {
			fmt.Fprintln(cli.stderr, "Error getting suggestions for", path, ":", err)
			os.Exit(1)
		}
		fmt.Fprint(cli.stdout, s)
	}
}
