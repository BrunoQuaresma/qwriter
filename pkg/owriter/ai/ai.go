package ai

import (
	"context"
	"math"

	"github.com/sashabaranov/go-openai"
)

type Client interface {
	SetPrompt(string)
	SendMessage(context.Context, string) (string, error)
}

type OpenAI struct {
	client *openai.Client
	prompt string
}

func NewOpenAI(authToken string) *OpenAI {
	return &OpenAI{
		client: openai.NewClient(authToken),
	}
}

func (d *OpenAI) SetPrompt(prompt string) {
	d.prompt = prompt
}

func (d *OpenAI) SendMessage(ctx context.Context, message string) (string, error) {
	resp, err := d.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			// Set the temperature to the minimum value for optimal accuracy.
			// Related issue: https://github.com/sashabaranov/go-openai/issues/9
			Temperature: math.SmallestNonzeroFloat32,
			Model:       openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: d.prompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
