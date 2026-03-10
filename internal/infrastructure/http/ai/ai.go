package aiinfra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/devathh/xcoder/internal/domain/ai"
)

type Request struct {
	Model    string        `json:"model"`
	Messages []*ai.Message `json:"messages"`
}

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type AiRepository struct {
	model  string
	client *http.Client
}

func New(model string, client *http.Client) *AiRepository {
	return &AiRepository{
		model:  model,
		client: client,
	}
}

func (ar *AiRepository) SendMsg(ctx context.Context, msgs []*ai.Message) (*ai.Message, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	reqBody := Request{
		Model:    ar.model,
		Messages: msgs,
	}

	jsonBytes, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:3264/api/chat/completions", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create req: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := ar.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do req: %w", err)
	}
	defer resp.Body.Close()

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}

	if len(response.Choices) < 1 {
		return nil, fmt.Errorf("invalid len of choices")
	}

	return &ai.Message{
		Content: response.Choices[0].Message.Content,
	}, nil
}
