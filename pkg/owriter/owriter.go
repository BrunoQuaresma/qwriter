package owriter

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/BrunoQuaresma/openwritter/pkg/owriter/ai"
)

type Writer interface {
	Suggestions(text string) ([]Suggestion, error)
	Apply(text string, suggestions []Suggestion) (string, error)
}

type Suggestion struct {
	Original string `json:"original"`
	Value    string `json:"value"`
}

type owriter struct {
	ai ai.Client
}

func New(ai ai.Client) Writer {
	return &owriter{
		ai: ai,
	}
}

const prompt = "You will be provided with text or code containing user-facing text. Your task is to correct any grammar and spelling errors and enhance the clarity and tone of the copy." +
	"Focus only on the text visible to the user, ignoring any code-related issues or comments." +
	"After making improvements, return a JSON array where each element includes the original text and its corresponding suggestion in this format: { original: \"...\", value: \"...\" }."

func (w *owriter) Suggestions(text string) ([]Suggestion, error) {
	ctx := context.Background()

	w.ai.SetPrompt(prompt)
	resp, err := w.ai.SendMessage(ctx, text)
	if err != nil {
		return []Suggestion{}, err
	}

	suggestions := []Suggestion{}
	err = json.Unmarshal([]byte(resp), &suggestions)
	if err != nil {
		return []Suggestion{}, err
	}

	return suggestions, nil
}

func (w *owriter) Apply(text string, suggestions []Suggestion) (string, error) {
	for _, s := range suggestions {
		text = strings.ReplaceAll(text, s.Original, s.Value)
	}

	return text, nil
}
