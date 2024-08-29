package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (cli *CLI) text() *cobra.Command {
	return &cobra.Command{
		Use:   "text [text]",
		Short: "Review text and return improved version.",
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
}
