package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/devathh/xcoder/internal/domain/ai"
	"github.com/devathh/xcoder/internal/domain/repository"
	aiinfra "github.com/devathh/xcoder/internal/infrastructure/http/ai"
)

type Service interface {
	Exec(ctx context.Context, prompt string)
}

type service struct {
	log        *slog.Logger
	aiRepo     *aiinfra.AiRepository
	repository repository.Repository
}

func New(
	log *slog.Logger,
	aiRepo *aiinfra.AiRepository,
	repository repository.Repository,
) Service {
	return &service{
		log:        log,
		aiRepo:     aiRepo,
		repository: repository,
	}
}

type Response struct {
	CommandType string `json:"command_type"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	Status      string `json:"status"`
}

func (s *service) Exec(ctx context.Context, prompt string) {
	bytes, _ := os.ReadFile("start.md")
	startPrompt := string(bytes)

	var msgs []*ai.Message
	msgs = append(msgs, &ai.Message{
		Role:    "system",
		Content: startPrompt,
	})

	msgs = append(msgs, &ai.Message{
		Role:    "user",
		Content: prompt,
	})

	for {
		resp, err := s.aiRepo.SendMsg(ctx, msgs)
		if err != nil {
			s.log.Error("failed to send msg into ai", slog.String("error", err.Error()))
			msgs = append(msgs, &ai.Message{
				Role:    "user",
				Content: fmt.Sprintf("failed to senf msg into ai: %s", err.Error()),
			})
		}

		var assistentResp Response
		if err := json.Unmarshal([]byte(resp.Content), &assistentResp); err != nil {
			s.log.Error("failed to unmarshal resp from ai", slog.String("error", err.Error()), slog.String("msg", resp.Content))
			msgs = append(msgs, &ai.Message{
				Role:    "user",
				Content: fmt.Sprintf("failed to senf msg into ai: %s", err.Error()),
			})
		}

		if assistentResp.CommandType == "done" {
			s.log.Info("status", slog.String("status", assistentResp.Status))
			return
		}

		msgs = append(msgs, s.execCommand(ctx, &assistentResp))
	}
}

func (s *service) execCommand(ctx context.Context, res *Response) *ai.Message {
	switch res.CommandType {
	case "create_file":
		s.log.Info("status", slog.String("status", res.Status))
		s.log.Info("creating file...", slog.String("filename", res.Name))
		err := s.repository.CreateFile(ctx, []byte(res.Body), res.Name)
		if err != nil {
			return &ai.Message{
				Role:    "user",
				Content: err.Error(),
			}
		}
	case "update_file":
		s.log.Info("status", slog.String("status", res.Status))
		s.log.Info("updating file...", slog.String("filename", res.Name))
		err := s.repository.UpdateFile(ctx, []byte(res.Body), res.Name)
		if err != nil {
			return &ai.Message{
				Role:    "user",
				Content: err.Error(),
			}
		}
	case "delete_file":
		s.log.Info("status", slog.String("status", res.Status))
		s.log.Info("deleting file...", slog.String("filename", res.Name))
		err := s.repository.DeleteFile(ctx, res.Name)
		if err != nil {
			return &ai.Message{
				Role:    "user",
				Content: err.Error(),
			}
		}
	case "read_file":
		s.log.Info("status", slog.String("status", res.Status))
		s.log.Info("reading file...", slog.String("filename", res.Name))
		bytes, err := s.repository.ReadFile(ctx, res.Name)
		if err != nil {
			return &ai.Message{
				Role:    "user",
				Content: err.Error(),
			}
		}

		return &ai.Message{
			Role:    "user",
			Content: string(bytes),
		}
	case "create_dir":
		s.log.Info("status", slog.String("status", res.Status))
		s.log.Info("creating dir...", slog.String("dir", res.Name))
		if err := s.repository.CreateDir(ctx, res.Name); err != nil {
			return &ai.Message{
				Role:    "user",
				Content: err.Error(),
			}
		}
	case "delete_dir":
		s.log.Info("status", slog.String("status", res.Status))
		s.log.Info("deleting dir...", slog.String("dir", res.Name))
		if err := s.repository.DeleteDir(ctx, res.Name); err != nil {
			return &ai.Message{
				Role:    "user",
				Content: err.Error(),
			}
		}
	case "read_dir":
		s.log.Info("status", slog.String("status", res.Status))
		s.log.Info("reading dir...", slog.String("dir", res.Name))
		rawDirs, err := s.repository.ReadDir(ctx, res.Name)
		if err != nil {
			return &ai.Message{
				Role:    "user",
				Content: err.Error(),
			}
		}

		dirs := ""
		for _, file := range rawDirs {
			dirs += fmt.Sprintf("%s ", file)
		}

		return &ai.Message{
			Role:    "user",
			Content: dirs,
		}
	}

	return &ai.Message{
		Role:    "user",
		Content: "no errors",
	}
}
