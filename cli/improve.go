package cli

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

func (cli *CLI) improve() *cobra.Command {
	return &cobra.Command{
		Use:   "improve",
		Short: "Improve text.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath != "" {
				config, err := NewConfig(configPath)
				if err != nil {
					return fmt.Errorf("failed to load configuration file: %w", err)
				}
				cli.config = config
			}

			err := cli.setProfile(profile)
			if err != nil {
				return err
			}

			txt := args[0]
			s, err := cli.writer.Improve(txt)
			if err != nil {
				return fmt.Errorf("failed to get suggestions for %s: %w", txt, err)
			}
			txt, err = cli.writer.Apply(txt, s)
			if err != nil {
				return fmt.Errorf("failed to apply suggestions for %s: %w", txt, err)
			}
			fmt.Fprintln(cli.stdout, txt)

			if copy {
				err = clipboard.WriteAll(txt)
				if err != nil {
					return fmt.Errorf("failed to copy text to clipboard: %w", err)
				}
			}

			return nil
		},
	}
}
