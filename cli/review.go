package cli

import (
	"fmt"

	"github.com/BrunoQuaresma/openwritter/pkg/owriter"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/gookit/goutil/fsutil"
	"github.com/spf13/cobra"
)

func (cli *CLI) review() *cobra.Command {
	return &cobra.Command{
		Use:   "review [glob pattern]",
		Short: "Review files and displays suggestions.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p := args[0]
			// The CLI only handle files and not directories so we use the WithFilesOnly
			// option to raise an error if a directory is passed.
			matches, err := doublestar.FilepathGlob(p, doublestar.WithFilesOnly())
			if err != nil {
				return fmt.Errorf("failed to match files for pattern %s: %w", p, err)
			}
			if matches == nil {
				return fmt.Errorf("no files found for pattern: %s", p)
			}

			var suggestionsByPath = make(map[string][]owriter.Suggestion)
			for _, p := range matches {
				text := fsutil.ReadFile(p)
				s, err := cli.writer.Suggestions(string(text))
				if err != nil {
					return fmt.Errorf("failed to get suggestions for %s: %w", p, err)
				}
				suggestionsByPath[p] = s
			}

			// Print the suggestions
			i := 0
			for p, suggestions := range suggestionsByPath {
				// Add one line break between files
				if i > 0 {
					fmt.Fprintln(cli.stdout)
				}
				for j, s := range suggestions {
					// Add one line break between suggestions
					if j > 0 {
						fmt.Fprintln(cli.stdout)
					}
					fmt.Fprintf(cli.stdout, "%d/%d for %s\n- %s\n+ %s\n", j+1, len(suggestions), p, s.Original, s.Value)
				}
				i++
			}

			return nil
		},
	}
}
