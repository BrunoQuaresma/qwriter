package cli_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/BrunoQuaresma/openwritter/cli"
	"github.com/BrunoQuaresma/openwritter/pkg/owriter"
	"github.com/gookit/goutil/fsutil"
	"github.com/stretchr/testify/require"
)

func TestFileMatching(t *testing.T) {
	t.Parallel()

	// Set up files by using the file path as the content. This helps identify
	// which files are being read and passed to the writer, as the writer only has
	// access to the content, not the file path.
	var files = []string{
		"text-file.txt",
		"go-file.go",
		"go-file-2.go",
		"site/script.js",
		"site/script-2.js",
		"site/index.html",
		"site/docs/tutorial.md",
		"site/docs/intro.md",
	}
	tmpDir := t.TempDir()
	filePath := func(path string) string {
		return fmt.Sprintf("%s/%s", tmpDir, path)
	}
	for _, relativePath := range files {
		absolutePath := filePath(relativePath)
		err := fsutil.WriteFile(absolutePath, absolutePath, fsutil.DefaultFilePerm)
		require.NoError(t, err, fmt.Sprintf("failed to write file %s", absolutePath))
	}

	tc := []struct {
		pattern string
		matches []string
	}{
		{
			pattern: filePath("site/*.js"),
			matches: []string{
				filePath("site/script.js"),
				filePath("site/script-2.js"),
			},
		},
		{
			pattern: filePath("site/*"),
			matches: []string{
				filePath("site/index.html"),
				filePath("site/script.js"),
				filePath("site/script-2.js"),
			},
		},
	}

	var out bytes.Buffer
	var writer mockWriter
	cli := cli.New(cli.Options{
		Writer: &writer,
		Stdout: &out, // Redirect stdout to a buffer to avoid printing to the terminal.
	})
	for _, c := range tc {
		t.Run(c.pattern, func(t *testing.T) {
			cli.Run([]string{"owriter", "--files", c.pattern})
			require.ElementsMatch(t, writer.analyzedContent, c.matches)
			writer.Reset()
		})
	}
}

type mockWriter struct {
	// analyzedContent is a list of all the content that was analyzed by the
	// writer. This helps to verify that the writer is reading the correct
	// content.
	analyzedContent []string
}

func (m *mockWriter) Suggestions(text string) ([]owriter.Suggestion, error) {
	m.analyzedContent = append(m.analyzedContent, text)
	return []owriter.Suggestion{}, nil
}

func (m *mockWriter) Apply(text string, suggestions []owriter.Suggestion) (string, error) {
	return text, nil
}

// Reset clears the analyzed content list. Helps to reset the state of the
// writer between tests.
func (m *mockWriter) Reset() {
	m.analyzedContent = []string{}
}
