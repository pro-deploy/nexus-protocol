package conversation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Service handles AI conversations
type Service struct {
	logger *zap.Logger
	// In production, add AI service client
}

// Conversation represents a conversation session
type Conversation struct {
	ID          string                `json:"id"`
	UserID      string                `json:"user_id"`
	Title       string                `json:"title,omitempty"`
	Status      string                `json:"status"` // active, completed, archived
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	MessageCount int                  `json:"message_count"`
	LastMessage *Message              `json:"last_message,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Message represents a single message in conversation
type Message struct {
	ID            string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	Role          string    `json:"role"` // user, assistant, system
	Content       string    `json:"content"`
	Timestamp     time.Time `json:"timestamp"`
	TokenCount    int       `json:"token_count,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// TypingStatus represents typing status
type TypingStatus struct {
	ConversationID string `json:"conversation_id"`
	UserID         string `json:"user_id"`
	IsTyping       bool   `json:"is_typing"`
	Timestamp      time.Time `json:"timestamp"`
}

// NewService creates a new conversation service
func NewService(logger *zap.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

// CreateConversation creates a new conversation
func (s *Service) CreateConversation(ctx context.Context, userID string, initialMessage string, metadata map[string]interface{}) (*Conversation, *Message, error) {
	conversationID := uuid.New().String()

	conversation := &Conversation{
		ID:          conversationID,
		UserID:      userID,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		MessageCount: 0,
		Metadata:    metadata,
	}

	s.logger.Info("Conversation created",
		zap.String("conversation_id", conversationID),
		zap.String("user_id", userID))

	// If initial message provided, add it
	var firstMessage *Message
	if initialMessage != "" {
		var err error
		firstMessage, err = s.SendMessage(ctx, conversationID, userID, initialMessage)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to add initial message: %w", err)
		}
		conversation.MessageCount = 1
		conversation.LastMessage = firstMessage
	}

	return conversation, firstMessage, nil
}

// SendMessage sends a message to conversation
func (s *Service) SendMessage(ctx context.Context, conversationID, userID, content string) (*Message, error) {
	if content == "" {
		return nil, fmt.Errorf("message content cannot be empty")
	}

	// Verify conversation ownership
	conversation, err := s.getConversation(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("conversation not found: %w", err)
	}

	if conversation.UserID != userID {
		return nil, fmt.Errorf("access denied: conversation belongs to different user")
	}

	// Create user message
	userMessage := &Message{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		Role:           "user",
		Content:        content,
		Timestamp:      time.Now(),
		TokenCount:     s.estimateTokenCount(content),
	}

	// Save user message
	if err := s.saveMessage(ctx, userMessage); err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}

	// Generate AI response
	aiResponse, err := s.generateAIResponse(ctx, conversationID, userID, content, conversation)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI response: %w", err)
	}

	// Save AI message
	if err := s.saveMessage(ctx, aiResponse); err != nil {
		return nil, fmt.Errorf("failed to save AI message: %w", err)
	}

	// Update conversation
	conversation.MessageCount += 2
	conversation.LastMessage = aiResponse
	conversation.UpdatedAt = time.Now()
	if err := s.updateConversation(ctx, conversation); err != nil {
		s.logger.Error("Failed to update conversation", zap.Error(err))
	}

	s.logger.Info("Message sent to conversation",
		zap.String("conversation_id", conversationID),
		zap.String("user_id", userID),
		zap.Int("message_length", len(content)))

	return aiResponse, nil
}

// GetConversation retrieves a conversation
func (s *Service) GetConversation(ctx context.Context, conversationID, userID string) (*Conversation, error) {
	conversation, err := s.getConversation(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	if conversation.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	return conversation, nil
}

// GetConversationHistory retrieves message history for a conversation
func (s *Service) GetConversationHistory(ctx context.Context, conversationID, userID string, limit, offset int) ([]*Message, error) {
	// Verify access
	conversation, err := s.getConversation(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	if conversation.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	// Get messages
	messages, err := s.getConversationMessages(ctx, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// UpdateTypingStatus updates typing status for a user in conversation
func (s *Service) UpdateTypingStatus(ctx context.Context, conversationID, userID string, isTyping bool) error {
	s.logger.Debug("Typing status updated",
		zap.String("conversation_id", conversationID),
		zap.String("user_id", userID),
		zap.Bool("is_typing", isTyping))

	// In production, broadcast to other participants via WebSocket
	// For now, just log

	return nil
}

// ArchiveConversation archives a conversation
func (s *Service) ArchiveConversation(ctx context.Context, conversationID, userID string) error {
	conversation, err := s.getConversation(ctx, conversationID)
	if err != nil {
		return err
	}

	if conversation.UserID != userID {
		return fmt.Errorf("access denied")
	}

	conversation.Status = "archived"
	conversation.UpdatedAt = time.Now()

	if err := s.updateConversation(ctx, conversation); err != nil {
		return err
	}

	s.logger.Info("Conversation archived",
		zap.String("conversation_id", conversationID),
		zap.String("user_id", userID))

	return nil
}

// ListUserConversations lists conversations for a user
func (s *Service) ListUserConversations(ctx context.Context, userID string, limit, offset int, includeArchived bool) ([]*Conversation, error) {
	// Mock implementation - in production would query database
	conversations := []*Conversation{
		{
			ID:          uuid.New().String(),
			UserID:      userID,
			Title:       "Обсуждение рецепта борща",
			Status:      "active",
			CreatedAt:   time.Now().Add(-time.Hour),
			UpdatedAt:   time.Now().Add(-5 * time.Minute),
			MessageCount: 12,
			LastMessage: &Message{
				Role:      "assistant",
				Content:   "Традиционный украинский борщ готовится с...",
				Timestamp: time.Now().Add(-5 * time.Minute),
			},
		},
		{
			ID:          uuid.New().String(),
			UserID:      userID,
			Title:       "Планирование путешествия",
			Status:      "active",
			CreatedAt:   time.Now().Add(-2 * time.Hour),
			UpdatedAt:   time.Now().Add(-10 * time.Minute),
			MessageCount: 8,
			LastMessage: &Message{
				Role:      "user",
				Content:   "Какие отели в Париже ты порекомендуешь?",
				Timestamp: time.Now().Add(-10 * time.Minute),
			},
		},
	}

	return conversations, nil
}

// generateAIResponse generates AI response for user message
func (s *Service) generateAIResponse(ctx context.Context, conversationID, userID, userMessage string, conversation *Conversation) (*Message, error) {
	// Get conversation history for context
	history, err := s.getConversationMessages(ctx, conversationID, 10, 0) // Last 10 messages
	if err != nil {
		s.logger.Warn("Failed to get conversation history", zap.Error(err))
	}

	// Build context from history
	var contextMessages []string
	for _, msg := range history {
		role := "Human"
		if msg.Role == "assistant" {
			role = "Assistant"
		}
		contextMessages = append(contextMessages, fmt.Sprintf("%s: %s", role, msg.Content))
	}

	contextStr := strings.Join(contextMessages, "\n")

	// Generate response based on user message and context
	response := s.generateMockResponse(userMessage, contextStr)

	aiMessage := &Message{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		Role:           "assistant",
		Content:        response,
		Timestamp:      time.Now(),
		TokenCount:     s.estimateTokenCount(response),
		Metadata: map[string]interface{}{
			"model": "gpt-4", // In production, actual model used
		},
	}

	return aiMessage, nil
}

// generateMockResponse generates a mock AI response (in production, call actual AI service)
func (s *Service) generateMockResponse(userMessage, context string) string {
	userMessage = strings.ToLower(userMessage)

	// Simple keyword-based responses (in production, use actual AI)
	if strings.Contains(userMessage, "борщ") {
		return "Борщ - это традиционное украинское блюдо, которое готовится из свеклы, капусты, картофеля, моркови и говядины. Хотите, я расскажу подробный рецепт?"
	}

	if strings.Contains(userMessage, "отель") || strings.Contains(userMessage, "путешествие") {
		return "Я могу помочь вам найти подходящий отель или спланировать путешествие. Расскажите подробнее о ваших предпочтениях: бюджет, дата, направление?"
	}

	if strings.Contains(userMessage, "рецепт") {
		return "Я знаю множество рецептов! От простых салатов до сложных десертов. Какой рецепт вас интересует?"
	}

	if strings.Contains(userMessage, "купить") || strings.Contains(userMessage, "магазин") {
		return "Я могу помочь найти товары и услуги. Что именно вы ищете? Укажите категорию или конкретный товар."
	}

	// Generic response
	return "Я понимаю ваш запрос. Расскажите подробнее, и я постараюсь помочь вам найти нужную информацию или выполнить задачу."
}

// estimateTokenCount estimates token count for a message (rough approximation)
func (s *Service) estimateTokenCount(content string) int {
	// Rough estimation: ~4 characters per token for English, ~2 for Russian
	if s.containsCyrillic(content) {
		return len(content) / 2
	}
	return len(content) / 4
}

// containsCyrillic checks if text contains Cyrillic characters
func (s *Service) containsCyrillic(text string) bool {
	for _, r := range text {
		if r >= '\u0400' && r <= '\u04FF' {
			return true
		}
	}
	return false
}

// Mock database operations (in production, implement with actual database)

func (s *Service) getConversation(ctx context.Context, conversationID string) (*Conversation, error) {
	// Mock implementation
	return &Conversation{
		ID:          conversationID,
		UserID:      "user-123", // Would get from database
		Status:      "active",
		CreatedAt:   time.Now().Add(-time.Hour),
		UpdatedAt:   time.Now(),
		MessageCount: 5,
	}, nil
}

func (s *Service) saveMessage(ctx context.Context, message *Message) error {
	// Mock implementation - in production would save to database
	s.logger.Debug("Message saved",
		zap.String("message_id", message.ID),
		zap.String("conversation_id", message.ConversationID),
		zap.String("role", message.Role))
	return nil
}

func (s *Service) updateConversation(ctx context.Context, conversation *Conversation) error {
	// Mock implementation
	return nil
}

func (s *Service) getConversationMessages(ctx context.Context, conversationID string, limit, offset int) ([]*Message, error) {
	// Mock implementation - return some sample messages
	return []*Message{
		{
			ID:             uuid.New().String(),
			ConversationID: conversationID,
			Role:           "user",
			Content:        "Расскажи рецепт борща",
			Timestamp:      time.Now().Add(-10 * time.Minute),
			TokenCount:     4,
		},
		{
			ID:             uuid.New().String(),
			ConversationID: conversationID,
			Role:           "assistant",
			Content:        "Борщ - традиционное блюдо украинской кухни...",
			Timestamp:      time.Now().Add(-9 * time.Minute),
			TokenCount:     25,
		},
	}, nil
}
