package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nexus-protocol/server/pkg/types"
)

// CreateConversationHandler handles conversation creation
func CreateConversationHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r.Context())

		var req struct {
			InitialMessage string                 `json:"initial_message,omitempty"`
			Metadata       map[string]interface{} `json:"metadata,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Invalid request format",
			}, http.StatusBadRequest)
			return
		}

		conversation, firstMessage, err := services.(*Services).Conversations.CreateConversation(r.Context(), userID, req.InitialMessage, req.Metadata)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to create conversation",
				Details: err.Error(),
			}, http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"conversation": conversation,
		}

		if firstMessage != nil {
			response["first_message"] = firstMessage
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

// GetConversationHandler handles conversation retrieval
func GetConversationHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		conversationID := vars["conversationId"]

		if conversationID == "" {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Conversation ID is required",
			}, http.StatusBadRequest)
			return
		}

		userID := getUserIDFromContext(r.Context())

		conversation, err := services.(*Services).Conversations.GetConversation(r.Context(), conversationID, userID)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "NOT_FOUND",
				Type:    "NOT_FOUND",
				Message: "Conversation not found",
			}, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"conversation": conversation,
		})
	}
}

// SendMessageHandler handles sending messages to conversations
func SendMessageHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		conversationID := vars["conversationId"]

		if conversationID == "" {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Conversation ID is required",
			}, http.StatusBadRequest)
			return
		}

		var req struct {
			Message string `json:"message"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Invalid request format",
			}, http.StatusBadRequest)
			return
		}

		if req.Message == "" {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Message content is required",
			}, http.StatusBadRequest)
			return
		}

		userID := getUserIDFromContext(r.Context())

		message, err := services.(*Services).Conversations.SendMessage(r.Context(), conversationID, userID, req.Message)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to send message",
				Details: err.Error(),
			}, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": message,
		})
	}
}

// GetConversationHistoryHandler handles conversation history retrieval
func GetConversationHistoryHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		conversationID := vars["conversationId"]

		if conversationID == "" {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Conversation ID is required",
			}, http.StatusBadRequest)
			return
		}

		userID := getUserIDFromContext(r.Context())

		// Parse pagination parameters
		limit := 50  // default
		offset := 0  // default

		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
				limit = parsedLimit
			}
		}

		if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
			if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
				offset = parsedOffset
			}
		}

		messages, err := services.(*Services).Conversations.GetConversationHistory(r.Context(), conversationID, userID, limit, offset)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to get conversation history",
			}, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"conversation_id": conversationID,
			"messages":        messages,
			"pagination": map[string]interface{}{
				"limit":  limit,
				"offset": offset,
				"count":  len(messages),
			},
		})
	}
}

// UpdateTypingStatusHandler handles typing status updates
func UpdateTypingStatusHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		conversationID := vars["conversationId"]

		if conversationID == "" {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Conversation ID is required",
			}, http.StatusBadRequest)
			return
		}

		var req struct {
			IsTyping bool `json:"is_typing"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Invalid request format",
			}, http.StatusBadRequest)
			return
		}

		userID := getUserIDFromContext(r.Context())

		if err := services.(*Services).Conversations.UpdateTypingStatus(r.Context(), conversationID, userID, req.IsTyping); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to update typing status",
			}, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"conversation_id": conversationID,
			"user_id":         userID,
			"is_typing":       req.IsTyping,
			"message":         "Typing status updated",
		})
	}
}

// ArchiveConversationHandler handles conversation archiving
func ArchiveConversationHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		conversationID := vars["conversationId"]

		if conversationID == "" {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Conversation ID is required",
			}, http.StatusBadRequest)
			return
		}

		userID := getUserIDFromContext(r.Context())

		if err := services.(*Services).Conversations.ArchiveConversation(r.Context(), conversationID, userID); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to archive conversation",
			}, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"conversation_id": conversationID,
			"status":          "archived",
			"message":         "Conversation archived successfully",
		})
	}
}

// ListUserConversationsHandler handles listing user conversations
func ListUserConversationsHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r.Context())

		// Parse query parameters
		limit := 20  // default
		offset := 0  // default
		includeArchived := false

		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 50 {
				limit = parsedLimit
			}
		}

		if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
			if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
				offset = parsedOffset
			}
		}

		if archivedStr := r.URL.Query().Get("include_archived"); archivedStr == "true" {
			includeArchived = true
		}

		conversations, err := services.(*Services).Conversations.ListUserConversations(r.Context(), userID, limit, offset, includeArchived)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to list conversations",
			}, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"conversations": conversations,
			"pagination": map[string]interface{}{
				"limit":  limit,
				"offset": offset,
				"count":  len(conversations),
			},
		})
	}
}
