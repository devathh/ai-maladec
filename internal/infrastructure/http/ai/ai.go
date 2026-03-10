package aiinfra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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
	url    string
	client *http.Client
}

// New создает репозиторий с явным URL и таймаутом
func New(model string, apiURL string, timeout time.Duration) *AiRepository {
	return &AiRepository{
		model: model,
		url:   apiURL,
		client: &http.Client{
			Timeout: timeout,
		},
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

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", ar.url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create req: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := ar.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do req: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}

	if len(response.Choices) < 1 {
		return nil, fmt.Errorf("invalid len of choices")
	}

	if response.Choices[0].Message.Content == "" {
		return nil, fmt.Errorf("invalid response from ai")
	}

	return ai.NewMessage(ai.RoleAssistant, response.Choices[0].Message.Content), nil
}
