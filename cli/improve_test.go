package cli_test

import (
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/BrunoQuaresma/openwritter/cli"
	"github.com/BrunoQuaresma/openwritter/cli/testutils"
	"github.com/BrunoQuaresma/openwritter/pkg/qwriter"
	"github.com/stretchr/testify/require"
)

func TestImprove_Output(t *testing.T) {
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
	cli.Run([]string{"improve", text})

	// Then return the improved text without errors
	require.Empty(t, stdError.String())
	require.Equal(t, "improved text\n", stdOut.String())
}

func TestImprove_NonExistentProfile(t *testing.T) {
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

	// When passing an invalid profile to the command
	cli.Run([]string{"improve", "improve this text", "-p", "nonexistent"})

	// Then return an error message
	require.NotEmpty(t, stdError)
	require.Contains(t, stdError.String(), "profile nonexistent not found")
}

func TestImprove_ProfileFromConfig(t *testing.T) {
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

	// Setup config
	tmp := t.TempDir()
	configPath := path.Join(tmp, ".qwriter.yaml")
	yaml := `
profiles:
  - name: "ui-copy"
    description: "generate and review text for the ui"
`
	os.WriteFile(configPath, []byte(yaml), 0644)

	// Pass the config path to the command and a profile from the config
	text := "improve this text"
	w.SetSuggestions(text, []qwriter.Suggestion{
		{Original: text, Value: "improved text"},
	})
	cli.Run([]string{"improve", text, "-c", configPath, "-p", "ui-copy"})

	// Then return the improved text without errors
	require.Empty(t, stdError.String())
	require.Equal(t, "improved text\n", stdOut.String())
}
