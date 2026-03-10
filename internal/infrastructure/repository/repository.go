package repositoryinfra

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type Repository struct {
}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) CreateFile(ctx context.Context, bytes []byte, filename string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(filename)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for file: %w", err)
		}
	}

	file, err := os.Create(filename)
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

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return bytes, nil
}

func (r *Repository) DeleteFile(ctx context.Context, filename string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return os.Remove(filename)
}

func (r *Repository) CreateDir(ctx context.Context, dir string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to mkdir: %w", err)
	}

	return nil
}

func (r *Repository) DeleteDir(ctx context.Context, dir string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("failed to remove dir: %w", err)
	}

	return nil
}

func (r *Repository) ReadDir(ctx context.Context, dir string) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	cleanPath := filepath.Clean(dir)

	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("the dir not found: %s", cleanPath)
		}
		return nil, fmt.Errorf("failed to get dir %s: %w", cleanPath, err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("this is a file, not a dir: %s", cleanPath)
	}

	entries, err := os.ReadDir(cleanPath)
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