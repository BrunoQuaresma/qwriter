package cli

import (
	"fmt"

	"github.com/BrunoQuaresma/openwritter/pkg/qwriter"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/gookit/goutil/fsutil"
	"github.com/spf13/cobra"
)

func (cli *CLI) files() *cobra.Command {
	return &cobra.Command{
		Use:   "files [pattern]",
		Short: "Review files to get improvement suggestions.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pattern := args[0]
			matches, err := doublestar.FilepathGlob(pattern, doublestar.WithFilesOnly())
			if err != nil {
				return fmt.Errorf("failed to match files for pattern %s: %w", pattern, err)
			}
			if matches == nil {
				return fmt.Errorf("no files found for pattern: %s", pattern)
			}

			// Get suggestions for each file
			var suggestionsByPath = make(map[string][]qwriter.Suggestion)
			for _, p := range matches {
				text := fsutil.ReadFile(p)
				s, err := cli.writer.Suggestions(string(text))
				if err != nil {
					return fmt.Errorf("failed to get suggestions for %s: %w", p, err)
				}
				suggestionsByPath[p] = s
			}

			// Print suggestions
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
