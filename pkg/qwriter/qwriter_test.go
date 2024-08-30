package qwriter_test

import (
	"flag"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/BrunoQuaresma/openwritter/pkg/qwriter"
	"github.com/BrunoQuaresma/openwritter/pkg/qwriter/ai"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update .golden files")

func TestFix(t *testing.T) {
	t.Parallel()

	syntaxes := []string{
		"react",
	}

	for _, syntax := range syntaxes {
		t.Run(syntax, func(t *testing.T) {
			t.Parallel()

			// Given: a code input with syntax errors.
			var (
				client = qwriter.New(qwriter.Options{
					AI:      ai.NewOpenAI(os.Getenv("QWRITER_OPENAI_KEY")),
					Profile: qwriter.DefaultProfile,
				})
				inputFile  = path.Join("testdata", fmt.Sprintf("%s.input", syntax))
				goldenFile = path.Join("testdata", fmt.Sprintf("%s.golden", syntax))
			)
			inputData, err := os.ReadFile(inputFile)
			if err != nil {
				t.Errorf("failed to read %s: %v", inputFile, err)
			}

			// When: fix the input code.
			input := string(inputData)
			suggestions, err := client.Suggestions(input)
			require.NoError(t, err, "failed to get suggestions")
			resp, err := client.Apply(input, suggestions)
			if err != nil {
				t.Errorf("failed to fix the code: %v", err)
			}

			// Then: the fix should match the expected code in the golden file.
			if *update {
				if err := os.WriteFile(goldenFile, []byte(resp), 0644); err != nil {
					t.Errorf("failed to write %s: %v", goldenFile, err)
				}
			}
			goldenData, err := os.ReadFile(goldenFile)
			if err != nil {
				t.Errorf("failed to read %s: %v", goldenFile, err)
			}
			require.Equal(t, string(goldenData), resp)
		})
	}
}
