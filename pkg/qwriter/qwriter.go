package qwriter

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/BrunoQuaresma/openwritter/pkg/qwriter/ai"
)

type Writer interface {
	SetProfile(Profile)
	Suggestions(text string) ([]Suggestion, error)
	Apply(text string, suggestions []Suggestion) (string, error)
}

type Suggestion struct {
	Original string `json:"original"`
	Value    string `json:"value"`
}

type qwriter struct {
	ai      ai.Client
	profile Profile
}

type Options struct {
	AI      ai.Client
	Profile Profile
}

func New(o Options) Writer {
	return &qwriter{
		ai:      o.AI,
		profile: o.Profile,
	}
}

// Adding this prompt after the profile prompt ensures that the suggestions are
// returned in the expected format.
const prompt = "You will be provided with text or code. Your task is to %s " +
	"After making improvements, return a JSON array where each element includes the original text and its corresponding suggestion in this format: { original: \"...\", value: \"...\" }."

func (w *qwriter) Suggestions(text string) ([]Suggestion, error) {
	ctx := context.Background()

	w.ai.SetPrompt(fmt.Sprintf(prompt, w.profile.Description))
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

func (w *qwriter) Apply(text string, suggestions []Suggestion) (string, error) {
	for _, s := range suggestions {
		text = strings.ReplaceAll(text, s.Original, s.Value)
	}

	return text, nil
}

func (w *qwriter) SetProfile(p Profile) {
	w.profile = p
}
