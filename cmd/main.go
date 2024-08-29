package main

import (
	"os"

	"github.com/BrunoQuaresma/openwritter/cli"
	"github.com/BrunoQuaresma/openwritter/pkg/qwriter"
	"github.com/BrunoQuaresma/openwritter/pkg/qwriter/ai"
)

func main() {
	w := qwriter.New(ai.NewOpenAI(os.Getenv("QWRITER_OPENAI_KEY")))
	cli := cli.New(cli.Options{
		Writer: w,
	})
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
