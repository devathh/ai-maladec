package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/devathh/xcoder/internal/application/services"
	aiinfra "github.com/devathh/xcoder/internal/infrastructure/http/ai"
	repositoryinfra "github.com/devathh/xcoder/internal/infrastructure/repository"
	"github.com/devathh/xcoder/pkg/log"
)

func main() {
	logH, _ := log.SetupHandler(os.Stdout, "dev")
	log := slog.New(logH)

	r := aiinfra.New("qwen3.5-397b-a17b", http.DefaultClient)
	repo := repositoryinfra.New()
	service := services.New(log, r, repo)

	fmt.Printf("What to do?: ")

	var prompt string
	fmt.Scanf("%s", &prompt)

	service.Exec(context.Background(), prompt)

}
