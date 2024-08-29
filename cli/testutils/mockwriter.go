package testutils

import (
	"strings"

	"github.com/BrunoQuaresma/openwritter/pkg/qwriter"
)

type MockWriter struct {
	// analyzedContent is a list of all the content that was analyzed by the
	// writer. This helps to verify that the writer is reading the correct
	// content.
	AnalyzedContent []string
	// suggestionsByText maps text to a list of suggestions. This is used to
	// predict the suggestions that the writer will return for specific file
	// content during tests, ensuring accurate and consistent test results.
	SuggestionsByText map[string][]qwriter.Suggestion
}

func (m *MockWriter) Suggestions(text string) ([]qwriter.Suggestion, error) {
	m.AnalyzedContent = append(m.AnalyzedContent, text)
	return m.SuggestionsByText[text], nil
}

func (m *MockWriter) Apply(text string, suggestions []qwriter.Suggestion) (string, error) {
	for _, s := range suggestions {
		text = strings.ReplaceAll(text, s.Original, s.Value)
	}

	return text, nil
}

// SetSuggestions is a utility function that sets the suggestions for a specific
// file content. Helps to predict the suggestions that the writer will return
// during tests.
func (m *MockWriter) SetSuggestions(text string, suggestions []qwriter.Suggestion) {
	if m.SuggestionsByText == nil {
		m.SuggestionsByText = make(map[string][]qwriter.Suggestion)
	}
	m.SuggestionsByText[text] = suggestions
}

// Reset clears the analyzed content list. Helps to reset the state of the
// writer between tests.
func (m *MockWriter) Reset() {
	m.AnalyzedContent = []string{}
}
