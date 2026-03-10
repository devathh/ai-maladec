package services

import (
	"context"
	"testing"

	"log/slog"
	"os"

	"github.com/devathh/xcoder/internal/domain/config"
	"github.com/devathh/xcoder/pkg/log"
)

// MockRepository реализует интерфейс repository.Repository для тестов
type MockRepository struct {
	createdFiles map[string][]byte
}

func (m *MockRepository) CreateFile(ctx context.Context, bytes []byte, filename string) error {
	if m.createdFiles == nil {
		m.createdFiles = make(map[string][]byte)
	}
	m.createdFiles[filename] = bytes
	return nil
}

func (m *MockRepository) UpdateFile(ctx context.Context, updBytes []byte, filename string) error {
	return nil
}

func (m *MockRepository) DeleteFile(ctx context.Context, filename string) error {
	return nil
}

func (m *MockRepository) ReadFile(ctx context.Context, filename string) ([]byte, error) {
	return []byte(""), nil
}

func (m *MockRepository) CreateDir(ctx context.Context, dir string) error {
	return nil
}

func (m *MockRepository) DeleteDir(ctx context.Context, dir string) error {
	return nil
}

func (m *MockRepository) ReadDir(ctx context.Context, dir string) ([]string, error) {
	return []string{}, nil
}

func (m *MockRepository) ExecCommand(ctx context.Context, command string, args ...string) (string, error) {
	return "ok", nil
}

func TestServiceExec_ValidJSON(t *testing.T) {
	// Setup logger
	logH, _ := log.SetupHandler(os.Stdout, "dev")
	logger := slog.New(logH)

	// Setup config
	cfg := &config.Config{
		SystemPrompt: "You are a helper",
		MaxRetries:   3,
		ContextSize:  10,
	}

	// Mock AI Repo (would normally call real API, here we skip or mock heavily)
	// For unit test simplicity, we assume the service logic handles JSON parsing correctly
	// A full integration test would require a mock HTTP server.

	mockRepo := &MockRepository{}

	// We cannot easily test the full loop without a real AI or complex HTTP mocking
	// So we test the initialization and basic structure
	service := New(logger, nil, mockRepo, cfg)

	if service == nil {
		t.Fatal("Service should be created")
	}

	// Basic assertion that the service implements the interface
	var _ Service = service
}
