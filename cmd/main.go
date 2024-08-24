package main

import (
	"os"

	"github.com/BrunoQuaresma/openwritter/cli"
	"github.com/BrunoQuaresma/openwritter/pkg/owriter"
	"github.com/BrunoQuaresma/openwritter/pkg/owriter/ai"
)

func main() {
	w := owriter.New(ai.NewOpenAI(os.Getenv("OPENAI_KEY")))
	cli := cli.New(cli.Options{
		Writer: w,
	})
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
