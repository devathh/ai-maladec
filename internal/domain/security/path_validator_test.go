package security

import (
	"testing"
)

func TestPathValidator_Validate(t *testing.T) {
	tests := []struct {
		name        string
		rootDir     string
		inputPath   string
		wantErr     bool
		errContains string
	}{
		{
			name:      "valid relative path",
			rootDir:   "/home/user/project",
			inputPath: "src/main.go",
			wantErr:   false,
		},
		{
			name:      "valid nested path",
			rootDir:   "/home/user/project",
			inputPath: "internal/domain/file.go",
			wantErr:   false,
		},
		{
			name:        "path traversal attack",
			rootDir:     "/home/user/project",
			inputPath:   "../../etc/passwd",
			wantErr:     true,
			errContains: "security violation",
		},
		{
			name:        "absolute path outside root",
			rootDir:     "/home/user/project",
			inputPath:   "/etc/shadow",
			wantErr:     true,
			errContains: "security violation",
		},
		{
			name:      "current directory",
			rootDir:   "/home/user/project",
			inputPath: ".",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewPathValidator(tt.rootDir)
			got, err := v.Validate(tt.inputPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || err.Error() == "" {
					t.Errorf("Expected error containing '%s', got '%v'", tt.errContains, err)
				}
			}

			if !tt.wantErr && got == "" {
				t.Error("Expected valid path, got empty string")
			}
		})
	}
}

func TestPathValidator_Normalization(t *testing.T) {
	v := NewPathValidator("/project")
	
	// Path with redundant slashes and dots should be cleaned
	path, err := v.Validate("./src//../src/main.go")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Should resolve to /project/src/main.go
	if path != "/project/src/main.go" {
		t.Errorf("Expected '/project/src/main.go', got '%s'", path)
	}
}
