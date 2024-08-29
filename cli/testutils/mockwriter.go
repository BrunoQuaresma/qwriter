package testutils

import (
	"strings"

	"github.com/BrunoQuaresma/openwritter/pkg/qwriter"
)

type MockWriter struct {
	// suggestionsByText maps text to a list of suggestions. This is used to
	// predict the suggestions that the writer will return for specific file
	// content during tests, ensuring accurate and consistent test results.
	SuggestionsByText map[string][]qwriter.Suggestion
}

func (m *MockWriter) Suggestions(text string) ([]qwriter.Suggestion, error) {
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
