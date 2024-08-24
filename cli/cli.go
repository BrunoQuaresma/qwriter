package cli

import (
	"errors"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			files, _ := cmd.Flags().GetString(string(FlagFiles))
			if files != "" {
				return cli.reviewFiles(files)
			}
			return nil
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

// Mute disables the output of the CLI. This function is primarily used for
// testing purposes.
func (cli *CLI) MuteOut() {
	cli.cmd.SetOut(io.Discard)
}

func (cli *CLI) reviewFiles(pattern string) error {
	// The CLI only handle files and not directories so we use the WithFilesOnly
	// option to raise an error if a directory is passed.
	matches, err := doublestar.FilepathGlob(pattern, doublestar.WithFilesOnly())
	if err != nil {
		return errors.New(pattern + " is an invalid pattern: " + err.Error())
	}
	if matches == nil {
		return errors.New("no files found for pattern: " + pattern)
	}

	var suggestionsByPath = make(map[string][]owriter.Suggestion)
	for _, path := range matches {
		s, err := cli.writer.Suggestions(path)
		if err != nil {
			return errors.New("failed to get suggestions for " + path + ": " + err.Error())
		}
		suggestionsByPath[path] = s
	}

	i := 0
	for path, suggestions := range suggestionsByPath {
		// Add one line break between files
		if i > 0 {
			fmt.Fprintln(cli.stdout)
		}
		for j, s := range suggestions {
			// Add one line break between suggestions
			if j > 0 {
				fmt.Fprintln(cli.stdout)
			}
			fmt.Fprintf(cli.stdout, "%d/%d for %s\n- %s\n+ %s\n", j+1, len(suggestions), path, s.Original, s.Value)
		}
		i++
	}

	return nil
}
