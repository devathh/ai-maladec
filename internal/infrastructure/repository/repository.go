package repositoryinfra

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devathh/xcoder/internal/domain/security"
)

type Repository struct {
	validator *security.PathValidator
	guard     *security.CommandGuard
}

func New(rootDir string, allowedCommands []string) *Repository {
	return &Repository{
		validator: security.NewPathValidator(rootDir),
		guard:     security.NewCommandGuard(allowedCommands),
	}
}

func (r *Repository) CreateFile(ctx context.Context, bytes []byte, filename string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	safePath, err := r.validator.Validate(filename)
	if err != nil {
		return fmt.Errorf("path validation failed: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(safePath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for file: %w", err)
		}
	}

	file, err := os.Create(safePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(bytes); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (r *Repository) UpdateFile(ctx context.Context, updBytes []byte, filename string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	safePath, err := r.validator.Validate(filename)
	if err != nil {
		return fmt.Errorf("path validation failed: %w", err)
	}

	file, err := os.OpenFile(safePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(updBytes); err != nil {
		return fmt.Errorf("failed to update body of file: %w", err)
	}

	return nil
}

func (r *Repository) ReadFile(ctx context.Context, filename string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	safePath, err := r.validator.Validate(filename)
	if err != nil {
		return nil, fmt.Errorf("path validation failed: %w", err)
	}

	bytes, err := os.ReadFile(safePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return bytes, nil
}

func (r *Repository) DeleteFile(ctx context.Context, filename string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	safePath, err := r.validator.Validate(filename)
	if err != nil {
		return fmt.Errorf("path validation failed: %w", err)
	}

	return os.Remove(safePath)
}

func (r *Repository) CreateDir(ctx context.Context, dir string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	safePath, err := r.validator.Validate(dir)
	if err != nil {
		return fmt.Errorf("path validation failed: %w", err)
	}

	if err := os.MkdirAll(safePath, 0755); err != nil {
		return fmt.Errorf("failed to mkdir: %w", err)
	}

	return nil
}

func (r *Repository) DeleteDir(ctx context.Context, dir string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	safePath, err := r.validator.Validate(dir)
	if err != nil {
		return fmt.Errorf("path validation failed: %w", err)
	}

	if err := os.RemoveAll(safePath); err != nil {
		return fmt.Errorf("failed to remove dir: %w", err)
	}

	return nil
}

func (r *Repository) ReadDir(ctx context.Context, dir string) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	safePath, err := r.validator.Validate(dir)
	if err != nil {
		return nil, fmt.Errorf("path validation failed: %w", err)
	}

	info, err := os.Stat(safePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("the dir not found: %s", safePath)
		}
		return nil, fmt.Errorf("failed to get dir %s: %w", safePath, err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("this is a file, not a dir: %s", safePath)
	}

	entries, err := os.ReadDir(safePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %w", err)
	}

	result := make([]string, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()

		if entry.IsDir() {
			name += "/"
		}

		result = append(result, name)
	}

	sort.Strings(result)

	return result, nil
}

func (r *Repository) ExecCommand(ctx context.Context, command string, args ...string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	// Проверка безопасности команды
	if err := r.guard.Validate(command, args); err != nil {
		return "", err
	}

	cmd := exec.CommandContext(ctx, command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %w, output: %s", err, string(output))
	}

	return strings.TrimSpace(string(output)), nil
}
