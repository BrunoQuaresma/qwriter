package cli_test

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"testing"

	"github.com/BrunoQuaresma/openwritter/cli"
	"github.com/BrunoQuaresma/openwritter/cli/testutils"
	"github.com/BrunoQuaresma/openwritter/pkg/qwriter"
	"github.com/gookit/goutil/fsutil"
	"github.com/stretchr/testify/require"
)

var (
	updateFilesOut = flag.Bool("update_files_out", false, "update files output golden files")
)

func TestFiles_Matching(t *testing.T) {
	t.Parallel()

	// Setup files with each file including its own path in the content. This
	// helps identify which files are being read and passed to the writer, as the
	// writer only has access to the content, not the file path.
	temp := t.TempDir()

	var f = []file{
		{
			Path:    filepath.Join(temp, "text-file.txt"),
			Content: filepath.Join(temp, "text-file.txt"),
		},
		{
			Path:    filepath.Join(temp, "go-file.go"),
			Content: filepath.Join(temp, "go-file.go"),
		},
		{
			Path:    filepath.Join(temp, "go-file-2.go"),
			Content: filepath.Join(temp, "go-file-2.go"),
		},
		{
			Path:    filepath.Join(temp, "site/script.js"),
			Content: filepath.Join(temp, "site/script.js"),
		},
		{
			Path:    filepath.Join(temp, "site/script-2.js"),
			Content: filepath.Join(temp, "site/script-2.js"),
		},
		{
			Path:    filepath.Join(temp, "site/index.html"),
			Content: filepath.Join(temp, "site/index.html"),
		},
		{
			Path:    filepath.Join(temp, "site/docs/tutorial.md"),
			Content: filepath.Join(temp, "site/docs/tutorial.md"),
		},
		{
			Path:    filepath.Join(temp, "site/docs/intro.md"),
			Content: filepath.Join(temp, "site/docs/intro.md"),
		},
	}
	err := setupFiles(f)
	require.NoError(t, err, "failed to setup files")

	// Define test cases for file matching. Each test case specifies a pattern to
	// match files and the expected files that should be matched. If the pattern
	// is invalid, the test case should have the error flag set to true.
	tc := []struct {
		pattern string
		matches []string
		error   bool
	}{
		{
			pattern: filepath.Join(temp, "site/*.js"),
			matches: []string{
				f[3].Path,
				f[4].Path,
			},
		},
		{
			pattern: filepath.Join(temp, "site/*"),
			matches: []string{
				f[3].Path,
				f[4].Path,
				f[5].Path,
			},
		},
		{
			pattern: filepath.Join(temp, "site/**/*"),
			matches: []string{
				f[3].Path,
				f[4].Path,
				f[5].Path,
				f[6].Path,
				f[7].Path,
			},
		},
		{
			pattern: filepath.Join(temp, "site/docs"),
			error:   true,
		},
		{
			pattern: filepath.Join(temp, "site/**"),
			error:   true,
		},
		{
			// This pattern will not match any files so it is expected to return an
			// error.
			pattern: filepath.Join(temp, "site/*.png"),
			error:   true,
		},
		{
			// This pattern is invalid because it uses the @! operator which is not
			// supported by glob.
			pattern: filepath.Join(temp, "@!"),
			error:   true,
		},
	}

	// Execute the command "qwriter --files <pattern>" for each test case. Verify
	// that the writer is correctly targeting the specified files or returning an
	// error.
	var (
		stdError bytes.Buffer
		w        testutils.MockWriter
	)
	cli := cli.New(cli.Options{
		Writer: &w,
		Stderr: &stdError,
		Stdout: io.Discard,
	})
	for _, c := range tc {
		t.Run(c.pattern, func(t *testing.T) {
			cli.Run([]string{"files", c.pattern})

			if c.error {
				require.NotEmpty(t, stdError.String(), "error should be present")
			} else {
				require.Empty(t, stdError.String(), "error should not be present")
				require.ElementsMatch(t, w.AnalyzedContent, c.matches)
			}

			w.Reset()
		})
	}
}

func TestFiles_Output(t *testing.T) {
	t.Parallel()

	// Setup test folder in a predictable path for testing with golden files. This
	// ensures that the file paths remain consistent, as they are used in the test
	// output for comparison.
	temp := "/tmp/TestReview_Output"
	err := fsutil.Mkdir(temp, fsutil.DefaultDirPerm)
	require.NoError(t, err, fmt.Sprintf("failed to create test directory: %s", temp))
	t.Cleanup(func() {
		fsutil.SafeRemoveAll(temp)
	})
	// Setup files
	var f = []file{
		{
			Path:    filepath.Join(temp, "docs/tutorial.md"),
			Content: "This is a tutorial sample. With some tutorial testing text.",
		},
		{
			Path:    filepath.Join(temp, "docs/intro.md"),
			Content: "I think it is ok for now to have it under UserAutocomplete. You brought good arguments.",
		},
	}
	err = setupFiles(f)
	require.NoError(t, err, "failed to setup files")

	// Set predictable suggestions for each file
	var w testutils.MockWriter
	w.SetSuggestions(f[0].Content, []qwriter.Suggestion{
		{
			Original: "This is a tutorial sample. With some tutorial testing text.",
			Value:    "This is a tutorial sample with some test text.",
		},
	})
	w.SetSuggestions(f[1].Content, []qwriter.Suggestion{
		{
			Original: "I think it is ok for now to have it under UserAutocomplete.",
			Value:    "I think it's fine for now to keep it under UserAutocomplete.",
		},
		{
			Original: "You brought good arguments.",
			Value:    "You made some good points.",
		},
	})

	// Execute the command
	var stdOut, stdErr bytes.Buffer
	cli := cli.New(cli.Options{
		Writer: &w,
		Stdout: &stdOut,
		Stderr: &stdErr,
	})
	cli.Run([]string{"files", filepath.Join(temp, "docs/*.md")})
	require.Empty(t, stdErr.String(), "error should not be present")

	// Update golden files
	goldenPath := "testdata/files_out.golden"
	if *updateFilesOut {
		fsutil.WriteFile(goldenPath, stdOut.String(), fsutil.DefaultFilePerm)
	}

	// Verify output
	golden := fsutil.ReadFile(goldenPath)
	require.Equal(t, string(golden), stdOut.String())
}

type file struct {
	Path    string
	Content string
}

func setupFiles(files []file) error {
	for _, f := range files {
		err := fsutil.WriteFile(f.Path, f.Content, fsutil.DefaultFilePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
