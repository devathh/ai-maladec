package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/devathh/xcoder/internal/application/services"
	"github.com/devathh/xcoder/internal/domain/config"
	aiinfra "github.com/devathh/xcoder/internal/infrastructure/http/ai"
	repositoryinfra "github.com/devathh/xcoder/internal/infrastructure/repository"
	"github.com/devathh/xcoder/pkg/log"
)

func main() {
	// 1. Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Инициализация логгера
	logH, err := log.SetupHandler(os.Stdout, cfg.Env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to setup logger: %v\n", err)
		os.Exit(1)
	}
	logger := slog.New(logH)

	// 3. Инициализация инфраструктуры
	aiRepo := aiinfra.New(cfg.AIModel, cfg.AIAPIURL, cfg.AITimeout)
	fsRepo := repositoryinfra.New(cfg.ProjectRoot, cfg.AllowedCommands)

	// 4. Инициализация сервиса
	service := services.New(logger, aiRepo, fsRepo, cfg)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("AI Coder Console initialized.")
	fmt.Printf("Model: %s | Timeout: %s\n", cfg.AIModel, cfg.AITimeout)
	fmt.Println("Enter your task description (Ctrl+C to exit):")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				logger.Error("Scan error", "error", err.Error())
			}
			break
		}

		prompt := strings.TrimSpace(scanner.Text())
		if prompt == "" {
			continue
		}

		logger.Info("Processing request", "prompt", prompt)
		
		// Выполняем команду только если промпт не пуст
		service.Exec(context.Background(), prompt)
	}
}