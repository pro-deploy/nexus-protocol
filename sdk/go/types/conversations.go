package types

// CreateConversationRequest представляет запрос создания беседы
type CreateConversationRequest struct {
	Title       string                 `json:"title,omitempty"`
	BotID       string                 `json:"bot_id,omitempty"`
	Context     map[string]interface{} `json:"context,omitempty"`
	SystemPrompt string                `json:"system_prompt,omitempty"`
}

// Conversation представляет беседу
type Conversation struct {
	ID           string   `json:"id"`
	UserID       string   `json:"user_id"`
	BotID        string   `json:"bot_id,omitempty"`
	Title        string   `json:"title,omitempty"`
	Status       string   `json:"status"`
	MessageCount int32    `json:"message_count"`
	CreatedAt    string   `json:"created_at"`
	LastActivity string   `json:"last_activity,omitempty"`
	Messages     []Message `json:"messages,omitempty"`
}

// Message представляет сообщение в беседе
type Message struct {
	ID             string                 `json:"id"`
	ConversationID string                 `json:"conversation_id"`
	SenderType     string                 `json:"sender_type"` // user, assistant, system
	Type           string                 `json:"type"`       // text, voice, image
	Content        string                 `json:"content"`
	Status         string                 `json:"status,omitempty"` // sent, delivered, read, failed
	CreatedAt      string                 `json:"created_at"`
	Intent         string                 `json:"intent,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// SendMessageRequest представляет запрос отправки сообщения
type SendMessageRequest struct {
	Content    string                 `json:"content"`
	MessageType string                `json:"message_type,omitempty"` // text, voice, image
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// MessageResponse представляет ответ на отправку сообщения
type MessageResponse struct {
	ConversationID    string   `json:"conversation_id"`
	UserMessage      *Message `json:"user_message,omitempty"`
	AIResponse       *Message `json:"ai_response,omitempty"`
	ConversationStatus string `json:"conversation_status,omitempty"`
	TotalMessages    int32    `json:"total_messages,omitempty"`
}

