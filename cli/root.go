package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (cli *CLI) root() *cobra.Command {
	return &cobra.Command{
		Use:   "owriter [text to improve]",
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
}
