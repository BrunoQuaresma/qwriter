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

// Flags
var (
	configPath string
	profile    string
	copy       bool
)

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
		config: nil,
	}
	cli.cmd = &cobra.Command{
		Use:   "qwriter",
		Short: "QWriter CLI is a tool to generate and write text using OpenAI's GPT-4 model.",
	}
	cli.cmd.SetOut(cli.stdout)
	cli.cmd.SetErr(cli.stderr)
	cli.cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to the QWriter CLI configuration file")
	cli.cmd.PersistentFlags().StringVarP(&profile, "profile", "p", qwriter.DefaultProfile.Name, "Profile to use for reviewing and generating text")
	cli.cmd.PersistentFlags().BoolVar(&copy, "copy", false, "Copy the reviwed and generated text to the clipboard")
	cli.cmd.AddCommand(cli.improve())

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

func (cli *CLI) setProfile(name string) error {
	if name == qwriter.DefaultProfile.Name {
		cli.writer.SetProfile(qwriter.DefaultProfile)
		return nil
	}

	var profiles []qwriter.Profile
	if cli.config != nil {
		profiles = cli.config.Profiles
	}

	var selected qwriter.Profile
	for _, p := range profiles {
		if p.Name == name {
			selected = p
			break
		}
	}

	//	A profile is required; display an error if it is not found.
	if selected.Name != "" {
		cli.writer.SetProfile(selected)
		return nil
	}

	return fmt.Errorf("profile %s not found", name)
}
