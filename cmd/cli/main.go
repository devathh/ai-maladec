package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/devathh/xcoder/internal/application/services"
	aiinfra "github.com/devathh/xcoder/internal/infrastructure/http/ai"
	repositoryinfra "github.com/devathh/xcoder/internal/infrastructure/repository"
	"github.com/devathh/xcoder/pkg/log"
)

func main() {
	logH, err := log.SetupHandler(os.Stdout, "dev")
	if err != nil {
		panic(err)
	}
	logger := slog.New(logH)

	// Инициализация зависимостей
	// Модель можно вынести в переменные окружения
	model := os.Getenv("AI_MODEL")
	if model == "" {
		model = "qwen3.5-397b-a17b"
	}

	aiRepo := aiinfra.New(model, &http.Client{})
	repo := repositoryinfra.New()
	service := services.New(logger, aiRepo, repo)

	ctx := context.Background()

	fmt.Println("=== XCoder CLI ===")
	fmt.Println("Введите команду или описание задачи (exit для выхода):")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">> ")
		if !scanner.Scan() {
			break
		}

		prompt := strings.TrimSpace(scanner.Text())
		if prompt == "exit" || prompt == "quit" {
			fmt.Println("Завершение работы...")
			break
		}

		if prompt == "" {
			continue
		}

		// Выполнение задачи через сервис
		service.Exec(ctx, prompt)
		fmt.Println("Готово. Жду следующую команду.")
	}

	if err := scanner.Err(); err != nil {
		logger.Error("Ошибка чтения ввода", slog.String("error", err.Error()))
	}
}