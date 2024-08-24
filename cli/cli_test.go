package cli_test

import (
	"bytes"
	"flag"
	"fmt"
	"testing"

	"github.com/BrunoQuaresma/openwritter/cli"
	"github.com/BrunoQuaresma/openwritter/pkg/owriter"
	"github.com/gookit/goutil/fsutil"
	"github.com/stretchr/testify/require"
)

func TestFilesMatching(t *testing.T) {
	t.Parallel()

	// Setup files. Each file's content includes its own path to help identify
	// which files are being read and passed to the writer, since the writer only
	// has access to the content, not the file path.
	var files = []file{
		{"text-file.txt", "text-file.txt"},
		{"go-file.go", "go-file.go"},
		{"go-file-2.go", "go-file-2.go"},
		{"site/script.js", "site/script.js"},
		{"site/script-2.js", "site/script-2.js"},
		{"site/index.html", "site/index.html"},
		{"site/docs/tutorial.md", "site/docs/tutorial.md"},
		{"site/docs/intro.md", "site/docs/intro.md"},
	}
	path := setupFiles(t, files)

	// Define test cases for file matching. Each test case specifies a pattern to
	// match files and the expected files that should be matched. If the pattern
	// is invalid, the test case should have the error flag set to true.
	tc := []struct {
		pattern string
		matches []string
		error   bool
	}{
		{
			pattern: path("site/*.js"),
			matches: []string{
				path("site/script.js"),
				path("site/script-2.js"),
			},
		},
		{
			pattern: path("site/*"),
			matches: []string{
				path("site/index.html"),
				path("site/script.js"),
				path("site/script-2.js"),
			},
		},
		{
			pattern: path("site/**/*"),
			matches: []string{
				path("site/index.html"),
				path("site/script.js"),
				path("site/script-2.js"),
				path("site/docs/tutorial.md"),
				path("site/docs/intro.md"),
			},
		},
		{
			pattern: path("site/docs"),
			error:   true,
		},
		{
			pattern: path("site/**"),
			error:   true,
		},
		{
			// This pattern will not match any files so it is expected to return an
			// error.
			pattern: path("site/*.png"),
			error:   true,
		},
		{
			// This pattern is invalid because it uses the @! operator which is not
			// supported by glob.
			pattern: path("@!"),
			error:   true,
		},
	}

	// Execute the command "owriter --files <pattern>" for each test case. Verify
	// that the writer is correctly targeting the specified files or returning an
	// error.
	var (
		err    bytes.Buffer
		writer mockWriter
	)
	cli := cli.New(cli.Options{
		Writer: &writer,
		Stderr: &err,
	})
	cli.MuteOut()
	for _, c := range tc {
		t.Run(c.pattern, func(t *testing.T) {
			cli.Run([]string{"owriter", "--files", c.pattern})

			if c.error {
				require.NotEmpty(t, err.String(), "error should be present")
			} else {
				require.Empty(t, err.String(), "error should not be present")
				require.ElementsMatch(t, writer.analyzedContent, c.matches)
			}

			writer.reset()
		})
	}
}

var (
	updateFilesOutput = flag.Bool("update-files-output", false, "update files output golden file")
)

func TestFilesOutput(t *testing.T) {
	t.Parallel()

	// Setup files
	var files = []file{
		{"docs/tutorial.md", "This is a tutorial sample. With some tutorial testing text."},
		{"docs/intro.md", "I think it is ok for now to have it under UserAutocomplete. You brought good arguments."},
	}
	path := setupFiles(t, files)

	// Set predictable suggestions for each file
	var writer mockWriter
	writer.setSuggestions(path("docs/tutorial.md"), []owriter.Suggestion{
		{
			Original: "This is a tutorial sample. With some tutorial testing text.",
			Value:    "This is a tutorial sample with some test text.",
		},
	})
	writer.setSuggestions(path("docs/intro.md"), []owriter.Suggestion{
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
	var out, err bytes.Buffer
	cli := cli.New(cli.Options{
		Writer: &writer,
		Stdout: &out,
		Stderr: &err,
	})
	cli.Run([]string{"owriter", "--files", path("docs/*.md")})
	require.Empty(t, err.String(), "error should not be present")

	// Update golden files
	goldenFilepath := "testdata/filesout.golden"
	if *updateFilesOutput {
		fsutil.WriteFile(goldenFilepath, out.String(), fsutil.DefaultFilePerm)
	}

	// Verify output
	golden := fsutil.ReadFile(goldenFilepath)
	require.Equal(t, string(golden), out.String())
}

type mockWriter struct {
	// analyzedContent is a list of all the content that was analyzed by the
	// writer. This helps to verify that the writer is reading the correct
	// content.
	analyzedContent []string
	// suggestionsByFile is a map of file paths to a list of suggestions. This
	// helps to predict the suggestions that the writer will return for each file
	// during tests.
	suggestionsByFilepath map[string][]owriter.Suggestion
}

func (m *mockWriter) Suggestions(text string) ([]owriter.Suggestion, error) {
	m.analyzedContent = append(m.analyzedContent, text)
	return m.suggestionsByFilepath[text], nil
}

func (m *mockWriter) Apply(text string, suggestions []owriter.Suggestion) (string, error) {
	return text, nil
}

// setSuggestions is a utility function that sets the suggestions for a specific
// file path. Helps to predict the suggestions that the writer will return
// during tests.
func (m *mockWriter) setSuggestions(path string, suggestions []owriter.Suggestion) {
	if m.suggestionsByFilepath == nil {
		m.suggestionsByFilepath = make(map[string][]owriter.Suggestion)
	}
	m.suggestionsByFilepath[path] = suggestions
}

// reset clears the analyzed content list. Helps to reset the state of the
// writer between tests.
func (m *mockWriter) reset() {
	m.analyzedContent = []string{}
}

type file struct {
	RelativePath string
	Content      string
}

// setupFiles is a utility function that creates files in a temporary folder for testing purposes.
// It returns a helper function that converts a relative file path to its corresponding absolute path
// within the temporary folder, ensuring accurate file operations during tests.
func setupFiles(t *testing.T, files []file) func(string) string {
	tmpDir := t.TempDir()
	path := func(path string) string {
		return fmt.Sprintf("%s/%s", tmpDir, path)
	}
	for _, f := range files {
		absolutePath := path(f.RelativePath)
		err := fsutil.WriteFile(absolutePath, f.Content, fsutil.DefaultFilePerm)
		require.NoError(t, err, fmt.Sprintf("failed to write file %s", absolutePath))
	}
	return path
}
