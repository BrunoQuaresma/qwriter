package main

import (
	"fmt"
	"os"

	"github.com/BrunoQuaresma/openwritter/cli"
	"github.com/BrunoQuaresma/openwritter/pkg/qwriter"
	"github.com/BrunoQuaresma/openwritter/pkg/qwriter/ai"
)

func main() {
	stdout := os.Stdout
	stderr := os.Stderr

	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		fmt.Fprintln(stderr, "Please, set the OPENAI_API_KEY environment variable with your OpenAI key.")
		os.Exit(1)
	}
	w := qwriter.New(qwriter.Options{
		AI: ai.NewOpenAI(os.Getenv("OPENAI_API_KEY")),
	})
	cli, err := cli.New(cli.Options{
		Writer: w,
		Stdout: stdout,
		Stderr: stderr,
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		os.Exit(1)
	}

	if err = cli.Execute(); err != nil {
		os.Exit(1)
	}
}
