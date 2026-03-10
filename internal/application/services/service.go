package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/devathh/xcoder/internal/domain/ai"
	"github.com/devathh/xcoder/internal/domain/config"
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
	cfg        *config.Config
}

func New(
	log *slog.Logger,
	aiRepo *aiinfra.AiRepository,
	repository repository.Repository,
	cfg *config.Config,
) Service {
	return &service{
		log:        log,
		aiRepo:     aiRepo,
		repository: repository,
		cfg:        cfg,
	}
}

type Response struct {
	CommandType string   `json:"command_type"`
	Name        string   `json:"name"`
	Body        string   `json:"body"`
	Status      string   `json:"status"`
	Command     string   `json:"command,omitempty"`
	Args        []string `json:"args,omitempty"`
}

func (s *service) Exec(ctx context.Context, prompt string) {
	var msgs []*ai.Message
	msgs = append(msgs, ai.NewMessage(ai.RoleSystem, s.cfg.SystemPrompt))
	msgs = append(msgs, ai.NewMessage(ai.RoleUser, prompt))

	retryCount := 0
	maxRetries := s.cfg.MaxRetries

	for {
		if retryCount >= maxRetries {
			s.log.Error("max retries exceeded", slog.Int("count", retryCount))
			fmt.Println("Error: max retries exceeded.")
			return
		}

		resp, err := s.aiRepo.SendMsg(ctx, msgs)
		if err != nil {
			s.log.Error("failed to send msg into ai", slog.String("error", err.Error()))
			retryCount++
			msgs = append(msgs, ai.NewMessage(ai.RoleUser, fmt.Sprintf("your answer is invalid: %s", err.Error())))
			continue
		}

		var assistentResp Response
		if err := json.Unmarshal([]byte(resp.Content), &assistentResp); err != nil {
			s.log.Error("failed to unmarshal resp from ai", slog.String("error", err.Error()), slog.String("msg", resp.Content))
			retryCount++
			msgs = append(msgs, ai.NewMessage(ai.RoleUser, fmt.Sprintf("failed to unmarshal response: %s. Please return valid JSON only, no markdown blocks.", err.Error())))
			continue
		}

		retryCount = 0

		if assistentResp.CommandType == "done" {
			s.log.Info("task completed", slog.String("status", assistentResp.Status))
			fmt.Println(assistentResp.Status)
			return
		}

		// Вывод статуса только если он отличается от предыдущего действия (упрощено до просто вывода статуса)
		// Убраны повторяющиеся логи и смайлики
		fmt.Printf("%s\n", assistentResp.Status)

		msgResult := s.execCommand(ctx, &assistentResp)

		if len(msgs) > s.cfg.ContextSize {
			msgs = append([]*ai.Message{msgs[0]}, msgs[len(msgs)-s.cfg.ContextSize+1:]...)
		}

		msgs = append(msgs, resp)
		msgs = append(msgs, msgResult)
	}
}

func (s *service) execCommand(ctx context.Context, res *Response) *ai.Message {
	switch res.CommandType {
	case "create_file":
		s.log.Debug("creating file", slog.String("filename", res.Name))
		if err := s.repository.CreateFile(ctx, []byte(res.Body), res.Name); err != nil {
			return ai.NewMessage(ai.RoleUser, err.Error())
		}

	case "update_file":
		s.log.Debug("updating file", slog.String("filename", res.Name))
		if err := s.repository.UpdateFile(ctx, []byte(res.Body), res.Name); err != nil {
			return ai.NewMessage(ai.RoleUser, err.Error())
		}

	case "delete_file":
		s.log.Debug("deleting file", slog.String("filename", res.Name))
		if err := s.repository.DeleteFile(ctx, res.Name); err != nil {
			return ai.NewMessage(ai.RoleUser, err.Error())
		}

	case "read_file":
		s.log.Debug("reading file", slog.String("filename", res.Name))
		bytes, err := s.repository.ReadFile(ctx, res.Name)
		if err != nil {
			return ai.NewMessage(ai.RoleUser, err.Error())
		}
		return ai.NewMessage(ai.RoleUser, string(bytes))

	case "create_dir":
		s.log.Debug("creating directory", slog.String("dir", res.Name))
		if err := s.repository.CreateDir(ctx, res.Name); err != nil {
			return ai.NewMessage(ai.RoleUser, err.Error())
		}

	case "delete_dir":
		s.log.Debug("deleting directory", slog.String("dir", res.Name))
		if err := s.repository.DeleteDir(ctx, res.Name); err != nil {
			return ai.NewMessage(ai.RoleUser, err.Error())
		}

	case "read_dir":
		s.log.Debug("reading directory", slog.String("dir", res.Name))
		rawDirs, err := s.repository.ReadDir(ctx, res.Name)
		if err != nil {
			return ai.NewMessage(ai.RoleUser, err.Error())
		}
		return ai.NewMessage(ai.RoleUser, strings.Join(rawDirs, " "))

	case "exec_command":
		cmd := res.Command
		if cmd == "" {
			cmd = res.Name
		}
		args := res.Args
		if len(args) == 0 && res.Body != "" {
			args = strings.Fields(res.Body)
		}

		s.log.Debug("executing system command", slog.String("command", cmd), slog.Any("args", args))
		output, err := s.repository.ExecCommand(ctx, cmd, args...)
		if err != nil {
			return ai.NewMessage(ai.RoleUser, fmt.Sprintf("command failed: %v, output: %s", err, output))
		}

		// Вывод результата команды только если он есть
		if output != "" {
			fmt.Println(output)
		}
		return ai.NewMessage(ai.RoleUser, output)

	default:
		s.log.Warn("unknown command type", slog.String("type", res.CommandType))
		return ai.NewMessage(ai.RoleUser, fmt.Sprintf("unknown command type: %s", res.CommandType))
	}

	return ai.NewMessage(ai.RoleUser, "ok")
}