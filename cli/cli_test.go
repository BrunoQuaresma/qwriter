package cli_test

import (
	"bytes"
	"testing"

	"github.com/BrunoQuaresma/openwritter/cli"
	"github.com/BrunoQuaresma/openwritter/cli/testutils"
	"github.com/BrunoQuaresma/openwritter/pkg/qwriter"
	"github.com/stretchr/testify/require"
)

func TestCLI_TextOutput(t *testing.T) {
	t.Parallel()

	var (
		stdError bytes.Buffer
		stdOut   bytes.Buffer
		w        testutils.MockWriter
	)
	cli, err := cli.New(cli.Options{
		Writer: &w,
		Stderr: &stdError,
		Stdout: &stdOut,
	})
	require.NoError(t, err)

	// When passing a text to the command
	text := "improve this text"
	w.SetSuggestions(text, []qwriter.Suggestion{
		{Original: text, Value: "improved text"},
	})
	cli.Run([]string{text})

	// Then return the improved text without errors
	require.Empty(t, stdError.String())
	require.Equal(t, "improved text\n", stdOut.String())
}
