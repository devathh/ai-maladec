package ai

// Роли участников диалога вынесены в константы для безопасности и типобезопасности
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Message представляет сообщение в чате с AI
type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

// NewMessage фабрика для создания сообщений
func NewMessage(role Role, content string) *Message {
	return &Message{
		Role:    role,
		Content: content,
	}
}
