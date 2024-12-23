package main

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type ChatGPTClient struct {
	*openai.Client `json:"-"` // Exclude from JSON
}

func InitChatGptClient(apiKey string) *ChatGPTClient {
	client := openai.NewClient(apiKey) // Create the OpenAI client
	return &ChatGPTClient{Client: client}
}

// SendChatGptRequest sends a request to ChatGPT with the given prompt and message
func (c *ChatGPTClient) SendChatGptRequest(userPrompt, msg string) (string, error) {
	req := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: userPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: msg,
		},
	}
	resp, err := c.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4o,
			Messages: req,
		},
	)
	if err != nil {
		return "", fmt.Errorf("sending request to ChatGPT failed: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}
func (c *ChatGPTClient) SendChatGptRequestWithHistory(userPrompt string, history []openai.ChatCompletionMessage) (string, []openai.ChatCompletionMessage, error) {
	history = append(history, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userPrompt,
	})

	// Send the request to the ChatGPT API
	resp, err := c.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4,
			Messages: history,
		},
	)
	if err != nil {
		return "", nil, fmt.Errorf("sending request to ChatGPT failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", nil, fmt.Errorf("no response received from ChatGPT")
	}
	history = append(history, resp.Choices[0].Message)

	return resp.Choices[0].Message.Content, history, nil
}

func InitChatGptHistory(globalPrompt string) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: globalPrompt,
		},
	}
}
