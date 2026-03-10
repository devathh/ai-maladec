package security

import (
	"testing"
)

func TestCommandGuard_Validate(t *testing.T) {
	allowed := []string{"go", "git", "ls", "cat", "mkdir", "rm"}
	guard := NewCommandGuard(allowed)

	tests := []struct {
		name        string
		command     string
		args        []string
		wantErr     bool
		errContains string
	}{
		{
			name:    "allowed command",
			command: "go",
			args:    []string{"build"},
			wantErr: false,
		},
		{
			name:    "allowed command with path",
			command: "/usr/bin/go",
			args:    []string{"version"},
			wantErr: false,
		},
		{
			name:        "forbidden command",
			command:     "curl",
			args:        []string{"http://evil.com"},
			wantErr:     true,
			errContains: "not in the allowed list",
		},
		{
			name:        "pipe injection",
			command:     "ls",
			args:        []string{"-la | rm -rf /"},
			wantErr:     true,
			errContains: "dangerous shell characters",
		},
		{
			name:        "semicolon injection",
			command:     "cat",
			args:        []string{"file.txt; wget evil.com"},
			wantErr:     true,
			errContains: "dangerous shell characters",
		},
		{
			name:        "backtick injection",
			command:     "echo",
			args:        []string{"`rm -rf /`"},
			wantErr:     true,
			errContains: "dangerous shell characters",
		},
		{
			name:        "shell script execution",
			command:     "sh",
			args:        []string{"-c", "rm -rf /"},
			wantErr:     true,
			errContains: "executing shell scripts via -c is forbidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := guard.Validate(tt.command, tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || err.Error() == "" {
					t.Errorf("Expected error containing '%s', got '%v'", tt.errContains, err)
				}
			}
		})
	}
}
