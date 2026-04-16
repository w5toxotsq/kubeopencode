package analyzer

import (
	"context"
	"fmt"

	"github.com/kubeopencode/kubeopencode/internal/k8s"
	"github.com/sashabaranov/go-openai"
)

const systemPrompt = `You are an expert Kubernetes administrator and developer.
Analyze the provided Kubernetes resource manifest and provide:
1. A brief summary of what the resource does
2. Potential issues or misconfigurations
3. Security recommendations
4. Performance optimization suggestions
Be concise and actionable.`

type Analyzer struct {
	client *openai.Client
	model  string
}

func New(apiKey string) *Analyzer {
	return &Analyzer{
		client: openai.NewClient(apiKey),
		// Switched to GPT-4o for better analysis quality — worth the extra cost for me
		model: openai.GPT4o,
	}
}

func NewWithModel(apiKey, model string) *Analyzer {
	return &Analyzer{
		client: openai.NewClient(apiKey),
		model:  model,
	}
}

func (a *Analyzer) Analyze(ctx context.Context, resource *k8s.Resource) (string, error) {
	resourceJSON, err := resource.ToJSON()
	if err != nil {
		return "", fmt.Errorf("failed to serialize resource: %w", err)
	}

	userMessage := fmt.Sprintf("Analyze this Kubernetes %s resource:\n\n```json\n%s\n```",
		resource.Kind, resourceJSON)

	resp, err := a.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: a.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: userMessage},
		},
		// Bumped from 1024 to 2048 — responses were getting cut off for larger resources
		MaxTokens:   2048,
		Temperature: 0.3,
	})
	if err != nil {
		return "", fmt.Errorf("openai request failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from openai")
	}

	return resp.Choices[0].Message.Content, nil
}
