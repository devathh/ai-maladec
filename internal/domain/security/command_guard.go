package security

import (
	"fmt"
	"path/filepath"
	"strings"
)

// CommandGuard отвечает за проверку безопасности исполняемых команд
type CommandGuard struct {
	allowedCommands map[string]bool
}

// NewCommandGuard создает защитник с белым списком разрешенных команд
func NewCommandGuard(allowedCommands []string) *CommandGuard {
	allowed := make(map[string]bool, len(allowedCommands))
	for _, cmd := range allowedCommands {
		allowed[cmd] = true
	}
	return &CommandGuard{
		allowedCommands: allowed,
	}
}

// Validate проверяет команду и аргументы на безопасность
func (g *CommandGuard) Validate(command string, args []string) error {
	// Нормализуем команду (убираем пути, если переданы полные пути к бинарникам)
	// Исправление: используем стандартный filepath.Base вместо самописной функции
	baseCmd := filepath.Base(command)

	if !g.allowedCommands[baseCmd] {
		return fmt.Errorf("security violation: command '%s' is not in the allowed list", baseCmd)
	}

	// Запрещаем опасные флаги и конструкции
	for _, arg := range args {
		if strings.Contains(arg, "|") || strings.Contains(arg, ";") || strings.Contains(arg, "&&") || strings.Contains(arg, "`") {
			return fmt.Errorf("security violation: dangerous shell characters detected in argument '%s'", arg)
		}
		if strings.HasPrefix(arg, "-") && strings.Contains(arg, "c") && baseCmd == "sh" {
			return fmt.Errorf("security violation: executing shell scripts via -c is forbidden")
		}
	}

	return nil
}
