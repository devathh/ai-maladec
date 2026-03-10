package repository

import "context"

type Repository interface {
	// File's
	CreateFile(ctx context.Context, bytes []byte, filename string) error
	UpdateFile(ctx context.Context, updBytes []byte, filename string) error
	DeleteFile(ctx context.Context, filename string) error
	ReadFile(ctx context.Context, filename string) ([]byte, error)

	// Dir's
	CreateDir(ctx context.Context, dir string) error
	DeleteDir(ctx context.Context, dir string) error
	ReadDir(ctx context.Context, dir string) ([]string, error)
}
