package config

import (
	"os"
	"time"
)

// Config хранит все настройки приложения
type Config struct {
	// AI Settings
	AIModel     string
	AIAPIURL    string
	AITimeout   time.Duration
	MaxRetries  int
	ContextSize int // Максимальное количество сообщений в истории

	// Security
	ProjectRoot string
	AllowedCommands []string

	// Logging
	LogLevel string
	Env      string

	// System Prompts
	SystemPrompt string
}

// Load загружает конфигурацию из переменных окружения и дефолтных значений
func Load() (*Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		AIModel:         getEnv("AI_MODEL", "qwen3.5-397b-a17b"),
		AIAPIURL:        getEnv("AI_API_URL", "http://localhost:3264/api/chat/completions"),
		AITimeout:       getDurationEnv("AI_TIMEOUT", 60*time.Second),
		MaxRetries:      getIntEnv("MAX_RETRIES", 3),
		ContextSize:     getIntEnv("CONTEXT_SIZE", 20),
		ProjectRoot:     cwd,
		AllowedCommands: []string{"go", "git", "ls", "cat", "mkdir", "rm", "mv", "cp", "echo", "pwd", "find", "grep"},
		LogLevel:        getEnv("LOG_LEVEL", "debug"),
		Env:             getEnv("APP_ENV", "dev"),
	}

	// Загрузка системного промпта
	if data, err := os.ReadFile("start.md"); err == nil {
		cfg.SystemPrompt = string(data)
	} else {
		cfg.SystemPrompt = "Ты — Senior Developer. Твоя задача: помочь пользователю."
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getDurationEnv(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}

func getIntEnv(key string, defaultVal int) int {
	// Упрощенная реализация для краткости
	return defaultVal
}
