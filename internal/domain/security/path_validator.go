package security

import (
	"fmt"
	"path/filepath"
	"strings"
)

// PathValidator отвечает за валидацию путей файлов для предотвращения Path Traversal
type PathValidator struct {
	rootDir string
}

// NewPathValidator создает валидатор с указанным корневым каталогом
func NewPathValidator(rootDir string) *PathValidator {
	// Очищаем и приводим к абсолютному пути корневую директорию
	absRoot, _ := filepath.Abs(rootDir)
	return &PathValidator{
		rootDir: absRoot,
	}
}

// Validate проверяет, что путь находится внутри разрешенной корневой директории
func (v *PathValidator) Validate(path string) (string, error) {
	// Очищаем путь от ../ и лишних слешей
	cleanPath := filepath.Clean(path)
	
	// Если путь относительный, делаем его абсолютным относительно rootDir
	if !filepath.IsAbs(cleanPath) {
		cleanPath = filepath.Join(v.rootDir, cleanPath)
	}
	
	// Приводим к абсолютному пути для сравнения
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Проверяем, начинается ли путь с корневой директории
	if !strings.HasPrefix(absPath, v.rootDir) {
		return "", fmt.Errorf("security violation: path '%s' is outside allowed root '%s'", path, v.rootDir)
	}

	return absPath, nil
}
